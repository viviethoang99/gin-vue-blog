import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';

import '../constants/theme.dart';
import '../services/tag_summary.dart';

@client
class TagsPage extends StatefulComponent {
  const TagsPage({required this.tags, super.key});

  final List<Map<String, dynamic>> tags;

  @override
  State<TagsPage> createState() => _TagsPageState();

  @css
  static List<StyleRule> get styles => [
    css('.tags-page').styles(
      width: 100.percent,
      minHeight: 100.vh,
      padding: .only(top: 112.px, right: 80.px, bottom: 70.px, left: 80.px),
      margin: .zero,
      boxSizing: .borderBox,
      color: textSecondaryColor,
      backgroundColor: Colors.transparent,
    ),
    css('.tags-container').styles(
      width: 100.percent,
      maxWidth: 720.px,
      margin: .only(left: .auto, right: .auto),
    ),
    css('.tags-header').styles(
      display: .flex,
      margin: .only(bottom: 34.px),
      justifyContent: .spaceBetween,
      alignItems: .end,
      gap: .all(24.px),
    ),
    css('.tags-header h1').styles(
      color: textPrimaryColor,
      fontSize: 2.rem,
      fontWeight: .w600,
      lineHeight: 1.2.em,
    ),
    css('.tags-header p').styles(
      margin: .only(top: 8.px),
      color: textMutedColor,
      fontSize: 0.92.rem,
      lineHeight: 1.5.em,
    ),
    css('.tag-sort-controls').styles(
      display: .inlineFlex,
      border: .all(style: .solid, color: borderColor, width: 1.px),
      radius: .all(.circular(7.px)),
      overflow: .hidden,
      flex: .none,
      backgroundColor: surfaceColor,
    ),
    css('.tag-sort-button').styles(
      minHeight: 34.px,
      padding: .symmetric(horizontal: 12.px),
      border: .none,
      cursor: .pointer,
      transition: Transition.combine([
        Transition('color', duration: 160.ms, curve: .easeOut),
        Transition('background-color', duration: 160.ms, curve: .easeOut),
      ]),
      color: textMutedColor,
      fontFamily: .inherit,
      fontSize: 0.8.rem,
      backgroundColor: Colors.transparent,
    ),
    css('.tag-sort-button + .tag-sort-button').styles(
      border: .only(
        left: .solid(color: borderColor, width: 1.px),
      ),
    ),
    css('.tag-sort-button:hover, .tag-sort-button:focus-visible').styles(
      color: textPrimaryColor,
      backgroundColor: postHoverColor,
    ),
    css('.tag-sort-button.active').styles(
      color: textPrimaryColor,
      fontWeight: .w600,
      backgroundColor: pageAccentColor,
    ),
    css('.tag-list').styles(
      display: .flex,
      flexWrap: .wrap,
      gap: .all(12.px),
    ),
    css('.tag-chip').styles(
      display: .inlineFlex,
      minHeight: 36.px,
      padding: .only(top: 7.px, right: 8.px, bottom: 7.px, left: 13.px),
      boxSizing: .borderBox,
      border: .all(style: .solid, color: borderColor, width: 1.px),
      radius: .all(.circular(18.px)),
      transition: Transition.combine([
        Transition('border-color', duration: 180.ms, curve: .easeOut),
        Transition('background-color', duration: 180.ms, curve: .easeOut),
      ]),
      alignItems: .center,
      gap: .all(8.px),
      color: primaryColor,
      fontSize: 0.86.rem,
      fontWeight: .w500,
      textDecoration: .none,
      backgroundColor: surfaceColor,
    ),
    css('.tag-chip:hover').styles(
      border: .all(style: .solid, color: primaryColor, width: 1.px),
      backgroundColor: postHoverColor,
    ),
    css('.tag-count').styles(
      display: .flex,
      height: 22.px,
      minWidth: 22.px,
      padding: .symmetric(horizontal: 6.px),
      boxSizing: .borderBox,
      radius: .all(.circular(11.px)),
      justifyContent: .center,
      alignItems: .center,
      color: textPrimaryColor,
      fontSize: 0.72.rem,
      fontWeight: .w600,
      lineHeight: 1.em,
      backgroundColor: pageAccentColor,
    ),
    css('.tags-empty').styles(color: textMutedColor, fontSize: 0.92.rem),
    css.media(MediaQuery.screen(maxWidth: 700.px), [
      css('.tags-page').styles(
        padding: .only(top: 112.px, right: 28.px, bottom: 52.px, left: 28.px),
      ),
      css('.tags-page .tags-header').styles(alignItems: .center, gap: .all(12.px)),
      css('.tags-page .tags-header h1').styles(fontSize: 1.75.rem),
      css('.tags-page .tag-sort-button').styles(
        minHeight: 32.px,
        padding: .symmetric(horizontal: 10.px),
      ),
      css('.tags-page .tag-list').styles(gap: .all(10.px)),
    ]),
    css.media(const MediaQuery.raw('(prefers-reduced-motion: reduce)'), [
      css('.tag-chip, .tag-sort-button').styles(raw: const {'transition': 'none'}),
    ]),
  ];
}

class _TagsPageState extends State<TagsPage> {
  TagSort _sort = TagSort.count;

  List<TagSummary> get _tags => [
    for (final tag in component.tags)
      TagSummary(
        id: tag['id'] as int,
        name: tag['name'] as String,
        count: tag['count'] as int,
      ),
  ];

  @override
  Component build(BuildContext context) {
    final tags = sortTagSummaries(_tags, _sort);

    return section(classes: 'tags-page', [
      main_(classes: 'tags-container', [
        header(classes: 'tags-header', [
          div([
            h1([.text('Tag')]),
            p([.text('${tags.length} tag từ các bài viết đã xuất bản')]),
          ]),
          div(
            classes: 'tag-sort-controls',
            attributes: const {'role': 'group', 'aria-label': 'Sắp xếp tag'},
            [
              _sortButton(TagSort.count, 'Số lượng'),
              _sortButton(TagSort.name, 'Tên'),
            ],
          ),
        ]),
        if (tags.isEmpty)
          p(classes: 'tags-empty', [.text('Chưa có tag nào.')])
        else
          div(classes: 'tag-list', [
            for (final tag in tags)
              a(href: '/tags/${tag.id}', classes: 'tag-chip', [
                span(classes: 'tag-name', [.text('#${tag.name}')]),
                span(classes: 'tag-count', [.text('${tag.count}')]),
              ]),
          ]),
      ]),
    ]);
  }

  Component _sortButton(TagSort value, String label) {
    final isActive = _sort == value;
    return button(
      type: .button,
      classes: 'tag-sort-button${isActive ? ' active' : ''}',
      attributes: {'aria-pressed': '$isActive'},
      onClick: () {
        if (_sort == value) return;
        setState(() => _sort = value);
      },
      [.text(label)],
    );
  }
}
