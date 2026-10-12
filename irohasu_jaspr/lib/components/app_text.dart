import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';

enum AppTextVariant { display, headline, title, body, label, caption }

enum AppTextTone { primary, secondary, accent, danger, success }

class AppText extends StatelessComponent {
  const AppText(
    this.value, {
    this.variant = AppTextVariant.body,
    this.tone = AppTextTone.primary,
    this.classes,
    super.key,
  });

  final String value;
  final AppTextVariant variant;
  final AppTextTone tone;
  final String? classes;

  @override
  Component build(BuildContext context) {
    final className = ['app-text', 'app-text--${variant.name}', 'app-text--${tone.name}', ?classes].join(' ');

    return switch (variant) {
      AppTextVariant.display => h1(classes: className, [.text(value)]),
      AppTextVariant.headline => h2(classes: className, [.text(value)]),
      AppTextVariant.title => h3(classes: className, [.text(value)]),
      AppTextVariant.body => p(classes: className, [.text(value)]),
      AppTextVariant.label || AppTextVariant.caption => span(classes: className, [.text(value)]),
    };
  }

  @css
  static List<StyleRule> get styles => [
    css('.app-text').styles(margin: .zero, color: const Color.variable('--color-text-primary')),
    css('.app-text--display').styles(fontSize: 3.5.rem, fontWeight: .w700, lineHeight: 1.1.em),
    css('.app-text--headline').styles(fontSize: 2.rem, fontWeight: .w700, lineHeight: 1.2.em),
    css('.app-text--title').styles(fontSize: 1.25.rem, fontWeight: .w700, lineHeight: 1.35.em),
    css('.app-text--body').styles(fontSize: 1.rem, lineHeight: 1.65.em),
    css('.app-text--label').styles(fontSize: 0.875.rem, fontWeight: .w700, lineHeight: 1.4.em),
    css('.app-text--caption').styles(fontSize: 0.75.rem, lineHeight: 1.4.em),
    css('.app-text--secondary').styles(color: const Color.variable('--color-text-secondary')),
    css('.app-text--accent').styles(color: const Color.variable('--color-primary')),
    css('.app-text--danger').styles(color: const Color.variable('--color-danger')),
    css('.app-text--success').styles(color: const Color.variable('--color-success')),
    css.media(MediaQuery.screen(maxWidth: 600.px), [
      css('.app-text--display').styles(fontSize: 2.5.rem),
      css('.app-text--headline').styles(fontSize: 1.65.rem),
    ]),
  ];
}
