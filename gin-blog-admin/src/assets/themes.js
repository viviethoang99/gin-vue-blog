import { darkTheme, lightTheme } from 'naive-ui'

function getNaiveThemeOverrides(darkMode) {
  const baseTheme = darkMode ? darkTheme : lightTheme
  const common = baseTheme.common

  return {
    common: {
      primaryColor: '#316C72FF',
      primaryColorHover: '#316C72E3',
      primaryColorPressed: '#2B4C59FF',
      primaryColorSuppl: '#316C72E3',

      infoColor: '#2080F0FF',
      infoColorHover: '#4098FCFF',
      infoColorPressed: '#1060C9FF',
      infoColorSuppl: '#4098FCFF',

      successColor: '#18A058FF',
      successColorHover: '#36AD6AFF',
      successColorPressed: '#0C7A43FF',
      successColorSuppl: '#36AD6AFF',

      warningColor: '#F0A020FF',
      warningColorHover: '#FCB040FF',
      warningColorPressed: '#C97C10FF',
      warningColorSuppl: '#FCB040FF',

      errorColor: '#D03050FF',
      errorColorHover: '#DE576DFF',
      errorColorPressed: '#AB1F3FFF',
      errorColorSuppl: '#DE576DFF',
    },
    Input: {
      color: common.cardColor,
      colorFocus: common.cardColor,
      textColor: common.textColor2,
      placeholderColor: common.placeholderColor,
      border: `1px solid ${common.borderColor}`,
      borderHover: `1px solid ${common.primaryColorHover}`,
      borderFocus: `1px solid ${common.primaryColorHover}`,
    },
  }
}

// TODO: 响应式
export default {
  header: {
    height: 60,
  },
  tags: {
    visible: true,
    height: 50,
  },
  naiveThemeOverrides: getNaiveThemeOverrides,
}
