import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';
import 'package:jaspr_router/jaspr_router.dart';

import '../config/environment.dart';
import '../constants/theme.dart';
import '../models/navigation_item.dart';
import 'article_search.dart';
import 'article_search_trigger.dart';
import 'theme_toggle.dart';

class Header extends StatelessComponent {
  const Header({
    required this.avatarUrl,
    required this.items,
    super.key,
  });

  final String avatarUrl;
  final List<NavigationItem> items;

  @override
  Component build(BuildContext context) {
    final activePath = context.url;

    return .fragment([
      aside(classes: 'site-sidebar', [
        div(classes: 'sidebar-inner', [
          if (avatarUrl.trim().isNotEmpty)
            Link(
              to: '/',
              child: img(src: avatarUrl, alt: 'Trang chủ', width: 34, height: 34, classes: 'sidebar-avatar'),
            ),
          nav(classes: 'sidebar-nav', _navigationLinks(activePath)),
          div(classes: 'sidebar-tools', [
            const ArticleSearchTrigger(),
            const ThemeToggle(),
          ]),
        ]),
      ]),
      details(classes: 'mobile-navigation', [
        summary(attributes: const {'aria-label': 'Mở menu điều hướng'}, [.text('☰')]),
        nav(_navigationLinks(activePath)),
        div(classes: 'mobile-navigation-tools', [
          const ArticleSearchTrigger(),
          const ThemeToggle(),
        ]),
      ]),
      ArticleSearchOverlay(apiBaseUrl: publicApiBaseUrl),
    ]);
  }

  List<Component> _navigationLinks(String activePath) => [
    for (final item in items)
      div(classes: 'navigation-item${_isActivePath(activePath, item.href) ? ' active' : ''}', [
        Link(
          to: item.href,
          attributes: {'aria-label': item.label},
          children: [
            span(classes: 'navigation-icon material-symbols-outlined', [.text(item.icon)]),
            span(classes: 'navigation-label', [.text(item.label)]),
          ],
        ),
      ]),
  ];

  bool _isActivePath(String activePath, String itemPath) {
    return activePath == itemPath || (itemPath != '/' && activePath.startsWith('$itemPath/'));
  }

  @css
  static List<StyleRule> get styles => [
    css('.site-sidebar').styles(
      display: .flex,
      position: .fixed(top: 0.px, bottom: 0.px, left: 0.px),
      zIndex: ZIndex(20),
      width: 96.px,
      pointerEvents: .none,
      justifyContent: .center,
      alignItems: .center,
    ),
    css('.sidebar-inner').styles(
      display: .flex,
      pointerEvents: .auto,
      flexDirection: .column,
      alignItems: .center,
      gap: .all(26.px),
    ),
    css('.sidebar-avatar').styles(
      display: .block,
      width: 34.px,
      height: 34.px,
      padding: .all(3.px),
      boxSizing: .borderBox,
      border: .all(style: .solid, color: borderColor, width: 1.px),
      radius: .all(.circular(50.percent)),
      backgroundColor: surfaceColor,
      raw: const {'object-fit': 'contain'},
    ),
    css('.sidebar-nav').styles(
      display: .flex,
      flexDirection: .column,
      alignItems: .center,
      gap: .all(8.px),
    ),
    css('.sidebar-nav .navigation-item').styles(position: .relative()),
    css('.sidebar-nav a').styles(
      display: .flex,
      width: 40.px,
      height: 40.px,
      radius: .all(.circular(6.px)),
      transition: Transition('background-color', duration: 180.ms, curve: .easeOut),
      justifyContent: .center,
      alignItems: .center,
      color: primaryColor,
      textDecoration: TextDecoration(line: .none),
      backgroundColor: Colors.transparent,
    ),
    css('.sidebar-nav a:hover, .sidebar-nav a:focus-visible').styles(backgroundColor: postHoverColor),
    css('.sidebar-nav .navigation-icon').styles(
      fontSize: 1.4.rem,
      lineHeight: 1.em,
      raw: const {
        'font-variation-settings': "'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24",
      },
    ),
    css('.sidebar-nav .navigation-label').styles(
      position: .absolute(left: 48.px),
      zIndex: ZIndex(1),
      padding: .symmetric(vertical: 7.px, horizontal: 10.px),
      border: .all(style: .solid, color: borderColor, width: 1.px),
      radius: .all(.circular(6.px)),
      opacity: 0,
      shadow: BoxShadow(offsetX: 0.px, offsetY: 8.px, blur: 20.px, color: const Color('#0000001f')),
      pointerEvents: .none,
      transition: Transition.combine([
        Transition('opacity', duration: 140.ms, curve: .easeOut),
        Transition('transform', duration: 140.ms, curve: .easeOut),
      ]),
      transform: .translate(x: (-4).px),
      color: textPrimaryColor,
      fontSize: 0.78.rem,
      fontWeight: .w500,
      lineHeight: 1.em,
      whiteSpace: .noWrap,
      backgroundColor: surfaceColor,
    ),
    css('.sidebar-nav a:hover .navigation-label, .sidebar-nav a:focus-visible .navigation-label').styles(
      opacity: 1,
      transform: .translate(x: 0.px),
    ),
    css('.sidebar-nav .active a').styles(color: textPrimaryColor, backgroundColor: pageAccentColor),
    css('.sidebar-tools').styles(
      display: .flex,
      flexDirection: .column,
      alignItems: .center,
      gap: .all(4.px),
    ),
    css('.sidebar-tools > .article-search-trigger, .site-sidebar .theme-toggle').styles(
      display: .flex,
      width: 40.px,
      height: 40.px,
      padding: .zero,
      border: .none,
      radius: .all(.circular(6.px)),
      cursor: .pointer,
      transition: Transition('background-color', duration: 180.ms, curve: .easeOut),
      justifyContent: .center,
      alignItems: .center,
      color: primaryColor,
      textDecoration: TextDecoration(line: .none),
      backgroundColor: Colors.transparent,
    ),
    css(
      '.sidebar-tools > .article-search-trigger:hover, .sidebar-tools > .article-search-trigger:focus-visible, .site-sidebar .theme-toggle:hover',
    ).styles(backgroundColor: postHoverColor),
    css('.sidebar-tools .material-symbols-outlined, .mobile-navigation-tools .material-symbols-outlined').styles(
      fontSize: 1.45.rem,
      lineHeight: 1.em,
      raw: const {
        'font-variation-settings': "'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24",
      },
    ),
    css('.mobile-navigation').styles(display: .none),
    css.media(MediaQuery.screen(maxWidth: 700.px), [
      css('.site-sidebar').styles(display: .none),
      css('.mobile-navigation').styles(
        display: .block,
        position: .fixed(top: 27.px, left: 28.px),
        zIndex: ZIndex(30),
        color: primaryColor,
      ),
      css('.mobile-navigation summary').styles(
        width: 28.px,
        height: 28.px,
        cursor: .pointer,
        userSelect: .none,
        listStyle: .none,
        color: primaryColor,
        fontSize: 1.1.rem,
        lineHeight: 28.px,
        raw: const {'list-style': 'none'},
      ),
      css('.mobile-navigation summary::-webkit-details-marker').styles(display: .none),
      css('.mobile-navigation nav').styles(
        display: .flex,
        width: 172.px,
        padding: .all(14.px),
        margin: .only(top: 10.px),
        border: .all(style: .solid, color: borderColor, width: 1.px),
        radius: .all(.circular(10.px)),
        shadow: BoxShadow(offsetX: 0.px, offsetY: 12.px, blur: 30.px, color: const Color('#00000033')),
        flexDirection: .column,
        gap: .all(5.px),
        backgroundColor: surfaceColor,
      ),
      css('.mobile-navigation-tools').styles(
        display: .flex,
        width: 172.px,
        padding: .only(right: 14.px, bottom: 12.px, left: 14.px),
        boxSizing: .borderBox,
        border: .all(style: .solid, color: borderColor, width: 1.px),
        radius: .only(bottomLeft: .circular(10.px), bottomRight: .circular(10.px)),
        justifyContent: .end,
        gap: .all(4.px),
        backgroundColor: surfaceColor,
        raw: const {'border-top': 'none'},
      ),
      css('.mobile-navigation nav:has(+ .mobile-navigation-tools)').styles(
        margin: .only(top: 10.px),
        radius: .only(topLeft: .circular(10.px), topRight: .circular(10.px)),
        raw: const {'margin-bottom': '0'},
      ),
      css('.mobile-navigation-tools .article-search-trigger, .mobile-navigation-tools .theme-toggle').styles(
        display: .flex,
        width: 36.px,
        height: 36.px,
        padding: .zero,
        border: .none,
        radius: .all(.circular(6.px)),
        cursor: .pointer,
        justifyContent: .center,
        alignItems: .center,
        color: primaryColor,
        backgroundColor: Colors.transparent,
      ),
      css(
        '.mobile-navigation-tools .article-search-trigger:hover, .mobile-navigation-tools .theme-toggle:hover',
      ).styles(
        backgroundColor: postHoverColor,
      ),
      css('.mobile-navigation nav a').styles(
        display: .block,
        width: .auto,
        height: .auto,
        padding: .symmetric(vertical: 8.px, horizontal: 10.px),
        radius: .all(.circular(6.px)),
        color: textSecondaryColor,
        fontSize: 0.88.rem,
        textDecoration: TextDecoration(line: .none),
        backgroundColor: Colors.transparent,
      ),
      css('.mobile-navigation .navigation-icon').styles(display: .none),
      css('.mobile-navigation .navigation-label').styles(
        position: .static,
        padding: .zero,
        border: .none,
        opacity: 1,
        shadow: BoxShadow.none,
        pointerEvents: .auto,
        transform: .none,
        color: .inherit,
        fontSize: .inherit,
        fontWeight: .w500,
        lineHeight: .inherit,
        backgroundColor: Colors.transparent,
        raw: const {'transition': 'none'},
      ),
      css('.mobile-navigation nav .active a').styles(color: primaryColor, backgroundColor: pageAccentColor),
    ]),
    css.media(const MediaQuery.raw('(prefers-reduced-motion: reduce)'), [
      css('.sidebar-nav a, .sidebar-nav .navigation-label').styles(raw: const {'transition': 'none'}),
    ]),
  ];
}
