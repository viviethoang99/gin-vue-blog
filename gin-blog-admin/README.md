This admin project is based on: [https://github.com/zclzone/vue-naive-admin](https://github.com/zclzone/vue-naive-admin). Thanks to the original author for the open-source work.

## Project Routing

- Backend routes: the backend returns a base menu array; the frontend assembles it into accessible routes.
- Frontend routes: load predefined frontend routes; use `meta.requireAuth` to determine whether authentication is required, and determine roles on the frontend.

## Differences from Vue Naive Admin

Principle: One problem doesn’t need too many solutions. This project keeps only the most common/necessary solutions. If you need more, please add them yourself.

Based on Vue Naive Admin, this project has been streamlined and integrated with a real backend. Main changes include:

Overall structure:
- Removed Mock: a real backend is provided, so Mock is unnecessary.
- Removed the build folder: many plugins were removed (all unplugin), so it’s not needed.
- Integrated real backend data and added backend-generated routes.

Plugins:
- Removed all unplugin series: `unplugin-auto-import`, `unplugin-icons`, `unplugin-vue-components`.
- Removed `vite-plugin-html`, `vite-plugin-mock`, `vite-plugin-svg-icons`: not used in this project.
- Removed Prettier, unify on ESLint.
- Removed `@commitlint/cli`, `@commitlint/config-conventional`: non-essential; keep it lean.
- Removed `lint-staged`, `husky`: this is a subproject in a monorepo, pre-commit checks are not required here.
- Removed `@unocss/preset-rem-to-px`: generally unnecessary to convert font sizes.
- Added `taze`: used to upgrade dependencies.

Reasons for removing unplugin series:
- These plugins are not required for business functionality, they mainly improve developer convenience.
- To reduce coupling and dependency on plugins, improving portability across environments.
- They can speed up solo development, but tend to be harder to maintain and less friendly to collaborators.
- They may also introduce strange issues.

UnoCSS (`uno.config.js`): aim for minimal presets
- Removed `presetAttributify`.
- Removed `shortcuts`.
- Removed `rules`.
- Use `@unocss/reset` instead of `reset.css`.
- Unify icon usage via UnoCSS with `presetIcons`.

## Code Style

Why not Prettier? See Antfu’s post: [Why I don’t use Prettier](https://antfu.me/posts/why-not-prettier-zh)

ESLint config: [https://github.com/antfu/eslint-config](https://github.com/antfu/eslint-config) to minimize configuration.
