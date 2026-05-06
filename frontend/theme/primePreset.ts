// theme/primePreset.ts
import Aura from '@primeuix/themes/aura'
import { definePreset, palette } from '@primeuix/themes'
// import type { ComponentsDesignTokens } from '@primeuix/themes/types'

const brand = {
  primary: palette('#00A499'),
  orange: palette('#ea7600'),
  grey: palette('#394049'),
  purple: palette('#8C4799'),
  blue: palette('#498BCA'),
  red: palette('#c8102e'),
  yellow: palette('#eaaa00'),
  hover: {
    primary: palette('#00837A'),
    orange: palette('#B95E00'),
    purple: palette('#713382'),
    blue: palette('#356a9e'),
    red: palette('#96081b'),
  },
}
/* Se comentó paleta con otros colores y componentes por problemas con la versión de primevue (4) */
/*
const components: ComponentsDesignTokens = {
  button: {
    // Variante "filled" (la que se usa sin variant)
    root: {
      primary: {
        background: '{primary.500}',
        hoverBackground: '{hover.primary.500}',
        borderColor: '{primary.500}',
        hoverBorderColor: '{hover.primary.500}',
        color: '{surface.0}',
      },
      secondary: {
        background: '{secondary.500}',
        hoverBackground: '{hover.secondary.500}',
        borderColor: '{secondary.500}',
        hoverBorderColor: '{hover.secondary.500}',
        color: '{surface.0}',
      },
      info: {
        background: '{secondary.500}',
        hoverBackground: '{hover.secondary.500}',
        borderColor: '{secondary.500}',
        hoverBorderColor: '{hover.secondary.500}',
        color: '{surface.0}',
      },
      danger: {
        background: '{red.500}',
        hoverBackground: '{hover.red.500}',
        borderColor: '{red.500}',
        hoverBorderColor: '{hover.red.500}',
        color: '{surface.0}',
      },
      warn: {
        background: '{orange.500}',
        hoverBackground: '{hover.orange.500}',
        borderColor: '{orange.500}',
        hoverBorderColor: '{hover.orange.500}',
        color: '{surface.0}',
      },
    },
    // (por si luego usas variant="outlined" o "text")
    outlined: {
      primary: {
        borderColor: '{primary.500}',
        color: '{primary.600}',
        hoverBackground: '{primary.50}',
      },
      secondary: {
        borderColor: '{secondary.500}',
        color: '{secondary.600}',
        hoverBackground: '{secondary.50}',
      },
    },
    text: {
      primary: {
        color: '{primary.600}',
        hoverBackground: '{primary.50}',
      },
      secondary: {
        color: '{secondary.600}',
        hoverBackground: '{secondary.50}',
      },
    },
  },
}
*/

const myPreset = definePreset(Aura, {
  /*
  primitive: {
    orange: brand.orange,
    grey: brand.grey,
    purple: brand.purple,
    blue: brand.blue,
    red: brand.red,
    yellow: brand.yellow,
    hover: {
      primary: brand.hover.primary,
      orange: brand.hover.orange,
      purple: brand.hover.purple,
      blue: brand.hover.blue,
      red: brand.hover.red,
    },
  },
  */
  semantic: { primary: brand.primary },
  /*
  components,
  */
})

export default myPreset
