import DefaultTheme from 'vitepress/theme'
import './custom.css'
import HeroInstall from './components/HeroInstall.vue'

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    app.component('HeroInstall', HeroInstall)
  }
}
