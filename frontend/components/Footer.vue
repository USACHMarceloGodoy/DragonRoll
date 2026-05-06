<template>
  <footer class="dr-footer">
    <p class="brand-line">
      <strong>DragonRoll</strong> · Compañero de mesa self-hosted
    </p>
    <p class="attribution">
      Esta aplicación incluye material del
      <a href="https://dnd.wizards.com/resources/systems-reference-document" target="_blank" rel="noopener">
        System Reference Document 5.1 (“SRD 5.1”)
      </a>
      por Wizards of the Coast LLC, bajo la
      <a href="https://creativecommons.org/licenses/by/4.0/legalcode" target="_blank" rel="noopener">
        Creative Commons Attribution 4.0 International License
      </a>.
      Compatible con la quinta edición.
    </p>
  </footer>
</template>

<style scoped>
.dr-footer {
  margin-top: 36px;
  padding: 16px 24px;
  background: #15110c;
  color: #8a8275;
  border-top: 1px solid #6b4a16;
  font-size: 11px;
  text-align: center;
}
.brand-line {
  margin: 0 0 6px;
  color: #c9c0ac;
  font-size: 12px;
}
.attribution {
  margin: 0;
  max-width: 720px;
  margin-left: auto; margin-right: auto;
  line-height: 1.5;
}
.attribution a { color: #b48a3c; text-decoration: underline; }
.attribution a:hover { color: #d4a04a; }
</style>
