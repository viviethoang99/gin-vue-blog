import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';

import '../constants/theme.dart';
import '../models/home_page_data.dart';

class SiteFooter extends StatelessComponent {
  const SiteFooter({required this.data, super.key});

  final SiteFooterData data;

  @override
  Component build(BuildContext context) {
    if (!_hasData) return .empty();

    final currentYear = DateTime.now().year;
    final startYear = data.startYear;
    final years = startYear == null
        ? null
        : startYear == currentYear
        ? '$currentYear'
        : '$startYear–$currentYear';

    return footer(classes: 'site-footer', [
      if (data.links.isNotEmpty)
        nav(classes: 'site-footer-links', [
          for (final link in data.links) a(href: link.href, [.text(link.label)]),
        ]),
      if (data.owner.isNotEmpty) div([.text('Copyright © ${years == null ? '' : '$years '}${data.owner}')]),
      div([.text('Built with Dart & Jaspr')]),
    ]);
  }

  bool get _hasData => data.owner.isNotEmpty || data.startYear != null || data.links.isNotEmpty;

  @css
  static List<StyleRule> get styles => [
    css('.site-footer').styles(
      width: 100.percent,
      padding: .only(right: 80.px, bottom: 24.px, left: 80.px),
      margin: .only(top: 72.px),
      boxSizing: .borderBox,
      flex: .none,
      color: textMutedColor,
      textAlign: .center,
      fontSize: 0.92.rem,
      lineHeight: 1.65.em,
    ),
    css('.site-footer-links').styles(
      display: .flex,
      justifyContent: .center,
      gap: .all(14.px),
    ),
    css('.site-footer-links a').styles(
      color: textMutedColor,
      textDecoration: .none,
    ),
    css('.site-footer-links a:hover, .site-footer-links a:focus-visible').styles(
      textDecoration: TextDecoration(line: .underline),
    ),
    css.media(MediaQuery.screen(maxWidth: 700.px), [
      css('.site-footer').styles(
        padding: .only(right: 28.px, bottom: 24.px, left: 28.px),
        margin: .only(top: 56.px),
      ),
    ]),
  ];
}
