#!/usr/bin/env python3
"""
parse_srd.py — extrae hechizos del SRD 5.1 oficial (CC-BY-4.0) a JSON
listo para el seeder de DragonRoll.

Uso:
    python parse_srd.py --input SRD5.1-CCBY4.0License-TT.json --output ../backend/data/seed/srd_spells.json

El SRD JSON (de github.com/Tabyltop/CC-SRD) viene plano: es una lista de
bloques type/level/subelements. Esta heurística:

  1. Busca el rango entre el h2 "Spell Descriptions" y el siguiente h2.
  2. Cada h4 dentro de ese rango es un nombre de hechizo.
  3. Los párrafos siguientes hasta el próximo h4 son su cuerpo.
  4. Se reconocen patrones "Casting Time:", "Range:", "Components:",
     "Duration:" y se extrae nivel/escuela del primer párrafo.
  5. "At Higher Levels." va a higher_level.

NOTA LEGAL: el SRD 5.1 está bajo CC-BY-4.0. El JSON de salida lleva
meta.source_label = 'SRD 5.1' para que el frontend pueda mostrar la
atribución requerida ("This work includes material taken from the SRD 5.1...").
"""
from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path
from typing import Iterable

LEVEL_WORDS = {
    "cantrip": 0,
    "1st-level": 1, "2nd-level": 2, "3rd-level": 3, "4th-level": 4,
    "5th-level": 5, "6th-level": 6, "7th-level": 7, "8th-level": 8,
    "9th-level": 9,
}

SCHOOLS_EN_ES = {
    "abjuration": "abjuración",
    "conjuration": "conjuración",
    "divination": "adivinación",
    "enchantment": "encantamiento",
    "evocation": "evocación",
    "illusion": "ilusión",
    "necromancy": "nigromancia",
    "transmutation": "transmutación",
}


def _text_of(block: dict) -> str:
    """Concatena el texto de todos los subelementos de un bloque."""
    return " ".join(
        s.get("text", "").strip() for s in block.get("subelements", [])
    ).strip()


def _find_block_range(blocks: list[dict], start_h2: str) -> tuple[int, int]:
    """Devuelve [inicio, fin) del rango entre el h2 dado y el siguiente h2."""
    start = None
    for i, b in enumerate(blocks):
        if b.get("type") == "h2" and _text_of(b) == start_h2:
            start = i
            break
    if start is None:
        raise ValueError(f"No se encontró h2 '{start_h2}'")
    end = len(blocks)
    for i in range(start + 1, len(blocks)):
        if blocks[i].get("type") == "h2":
            end = i
            break
    return start, end


def _split_into_spells(spell_blocks: list[dict]) -> Iterable[tuple[str, list[dict]]]:
    """
    Itera (nombre, [bloques del cuerpo]) por cada hechizo. Cada bloque
    h4 inicia un hechizo nuevo.
    """
    name = None
    body: list[dict] = []
    for b in spell_blocks:
        if b.get("type") == "h4":
            if name is not None:
                yield name, body
            name = _text_of(b).strip()
            body = []
        else:
            if name is not None:
                body.append(b)
    if name is not None:
        yield name, body


def _parse_first_paragraph(text: str) -> tuple[int, str, bool]:
    """
    El primer párrafo de un hechizo tiene la forma:
      "<level>-level <school>" o "<school> cantrip"
    A veces sufre un "(ritual)" o termina ahí.
    Retorna: (level, school_es, is_ritual).
    """
    t = text.lower().strip()
    is_ritual = "(ritual)" in t
    t = t.replace("(ritual)", "").strip()

    # cantrip case: "evocation cantrip"
    m = re.match(r"([a-z]+)\s+cantrip", t)
    if m:
        school = SCHOOLS_EN_ES.get(m.group(1), m.group(1))
        return 0, school, False

    # nivel: "2nd-level evocation"
    m = re.match(r"(\d+)(?:st|nd|rd|th)-level\s+([a-z]+)", t)
    if m:
        level = int(m.group(1))
        school = SCHOOLS_EN_ES.get(m.group(2), m.group(2))
        return level, school, is_ritual

    # fallback: extraer los enteros y palabras conocidas como podamos
    return 0, "evocación", is_ritual


def _extract_field(line: str, prefix: str) -> str | None:
    """Detecta líneas tipo 'Casting Time: 1 action' y devuelve el valor."""
    p = prefix.lower() + ":"
    if line.lower().startswith(p):
        return line[len(p):].strip()
    return None


def _parse_spell(name: str, body: list[dict]) -> dict | None:
    """
    Construye el dict del hechizo a partir de los bloques de su cuerpo.
    Devuelve None si parece corrupto (sin nivel ni componentes).
    """
    # Texto plano de los párrafos.
    paragraphs = [_text_of(b) for b in body if b.get("type") == "paragraph"]
    if not paragraphs:
        return None

    level, school, ritual_meta = _parse_first_paragraph(paragraphs[0])

    casting_time = "action"
    rng = ""
    components = ""
    duration = ""
    is_ritual = ritual_meta
    is_concentration = False
    higher_level = ""

    body_lines: list[str] = []
    for p in paragraphs[1:]:
        # Algunas conversiones del PDF concatenan "Duration: Instantaneous"
        # con el siguiente párrafo descriptivo. Detectamos los prefijos al
        # principio del texto.
        ct = _extract_field(p, "Casting Time")
        rg = _extract_field(p, "Range")
        cm = _extract_field(p, "Components")
        du = _extract_field(p, "Duration")
        if ct is not None:
            casting_time = ct
            continue
        if rg is not None:
            rng = rg
            continue
        if cm is not None:
            components = cm
            continue
        if du is not None:
            # Duration suele venir pegado a la descripción ("Instantaneous A
            # shimmering green arrow…"). Detectamos el patrón final típico:
            #   - "Instantaneous"
            #   - "Concentration, up to X minute(s)/hour(s)"
            #   - "X minute(s) / hour(s) / round(s) / day(s)"
            #   - "Until dispelled"
            #   - "Special"
            duration_re = re.compile(
                r"^(?:"
                r"Instantaneous"
                r"|Until dispelled(?: or triggered)?"
                r"|Special"
                r"|Concentration,?\s+up to\s+\d+\s+\w+"
                r"|Up to\s+\d+\s+(?:round|minute|hour|day)s?"
                r"|\d+\s+(?:round|minute|hour|day)s?"
                r")",
                re.IGNORECASE,
            )
            m = duration_re.match(du)
            if m:
                duration = m.group(0).strip()
                tail = du[m.end():].strip()
                if tail:
                    body_lines.append(tail)
            else:
                # fallback: cortamos en el primer punto seguido de mayúscula.
                split = re.split(r"(?<=[a-z])(?=[A-Z])", du, maxsplit=1)
                duration = split[0].strip()
                if len(split) > 1:
                    body_lines.append(split[1].strip())
            if "concentration" in duration.lower():
                is_concentration = True
            continue

        # "At Higher Levels." → higher_level
        if p.lower().startswith("at higher levels"):
            higher_level = p.split(".", 1)[1].strip() if "." in p else p
            continue

        body_lines.append(p)

    description = "\n\n".join(body_lines).strip()
    if not description:
        return None

    return {
        "tag": f"srd:{re.sub(r'[^a-z0-9]+', '-', name.lower()).strip('-')}",
        "name": name,
        "level": level,
        "school": school,
        "casting_time": casting_time,
        "range": rng,
        "components": components,
        "duration": duration,
        "description": description,
        "higher_level": higher_level,
        "ritual": is_ritual,
        "concentration": is_concentration,
    }


def parse_srd_spells(srd_path: Path) -> list[dict]:
    blocks = json.loads(srd_path.read_text(encoding="utf-8"))
    start, end = _find_block_range(blocks, "Spell Descriptions")
    spell_blocks = blocks[start + 1: end]
    spells: list[dict] = []
    for name, body in _split_into_spells(spell_blocks):
        spell = _parse_spell(name, body)
        if spell:
            spells.append(spell)
    return spells


def main() -> int:
    parser = argparse.ArgumentParser(description="Extrae hechizos del SRD 5.1 a JSON seed.")
    parser.add_argument("--input", required=True, type=Path, help="Ruta al SRD JSON oficial (Tabyltop/CC-SRD).")
    parser.add_argument("--output", required=True, type=Path, help="Ruta del JSON resultante.")
    args = parser.parse_args()

    if not args.input.exists():
        print(f"Error: no existe {args.input}", file=sys.stderr)
        return 1

    spells = parse_srd_spells(args.input)
    print(f"Parseados {len(spells)} hechizos.")

    # Forma compatible con el seeder homebrew: features/spells/races/classes.
    out = {
        "_attribution": (
            "This work includes material taken from the System Reference Document 5.1 "
            "(\"SRD 5.1\") by Wizards of the Coast LLC and available at "
            "https://dnd.wizards.com/resources/systems-reference-document. The SRD 5.1 "
            "is licensed under the Creative Commons Attribution 4.0 International "
            "License available at https://creativecommons.org/licenses/by/4.0/legalcode."
        ),
        "features": [],
        "spells": spells,
        "races": [],
        "classes": [],
    }
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(out, indent=2, ensure_ascii=False), encoding="utf-8")
    print(f"Escrito {args.output} ({args.output.stat().st_size // 1024} KB)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
