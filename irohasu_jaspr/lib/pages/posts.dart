import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';

import '../constants/theme.dart';
import '../models/home_page_data.dart';
import '../services/pagination.dart';

@client
class PostsPage extends StatefulComponent {
  const PostsPage({
    required this.posts,
    required this.authorName,
    required this.tagline,
    required this.avatarUrl,
    this.pageSize = 5,
    super.key,
  }) : assert(pageSize > 0);

  final List<Map<String, dynamic>> posts;
  final String authorName;
  final String tagline;
  final String avatarUrl;
  final int pageSize;

  @override
  State<PostsPage> createState() => _PostsPageState();

  @css
  static List<StyleRule> get styles => [
    css('.posts-page', [
      css('&').styles(
        width: 100.percent,
        minHeight: 100.vh,
        padding: .only(top: 112.px, right: 80.px, bottom: 70.px, left: 80.px),
        margin: .zero,
        boxSizing: .borderBox,
        color: textSecondaryColor,
        backgroundColor: Colors.transparent,
      ),
      css('a', [
        css('&').styles(
          color: primaryColor,
          textDecoration: TextDecoration(line: .none),
        ),
        css('&:hover').styles(textDecoration: TextDecoration(line: .underline)),
      ]),
      css('.posts-container').styles(
        width: 100.percent,
        maxWidth: 720.px,
        margin: .only(left: .auto, right: .auto),
      ),
      css('.posts-profile').styles(margin: .only(bottom: 34.px)),
      css('.posts-avatar').styles(
        display: .block,
        width: 100.px,
        height: 100.px,
        padding: .all(8.px),
        boxSizing: .borderBox,
        border: .all(style: .solid, color: borderColor, width: 1.px),
        radius: .all(.circular(50.percent)),
        backgroundColor: surfaceColor,
        raw: const {'object-fit': 'contain'},
      ),
      css('.posts-author-name').styles(
        margin: .only(top: 12.px),
        color: textPrimaryColor,
        fontFamily: const .list([FontFamily('Caveat'), FontFamilies.cursive]),
        fontSize: 2.4.rem,
        fontWeight: .w600,
        letterSpacing: 1.px,
        lineHeight: 1.em,
      ),
      css('.posts-tagline').styles(
        margin: .only(top: 7.px),
        color: textMutedColor,
        fontSize: 0.96.rem,
        lineHeight: 1.5.em,
      ),
      css('.archive-list').styles(
        display: .flex,
        flexDirection: .column,
        gap: .all(36.px),
      ),
      css('.empty-posts-message').styles(
        margin: .zero,
        color: textMutedColor,
        fontSize: 0.9.rem,
        lineHeight: 1.55.em,
      ),
      css('.archive-post h2').styles(
        margin: .only(bottom: 6.px),
        color: textPrimaryColor,
        fontSize: 1.12.rem,
        fontWeight: .w600,
        lineHeight: 1.35.em,
      ),
      css('.archive-post h2 a').styles(color: textPrimaryColor),
      css('.archive-excerpt').styles(
        margin: .zero,
        color: textSecondaryColor,
        fontSize: 0.9.rem,
        lineHeight: 1.55.em,
      ),
      css('.archive-meta').styles(
        display: .flex,
        margin: .only(top: 6.px),
        flexWrap: .wrap,
        gap: Gap(row: 3.px, column: 7.px),
        color: primaryColor,
        fontSize: 0.76.rem,
        lineHeight: 1.4.em,
      ),
      css('.archive-meta-item:not(:last-child)::after').styles(
        content: '·',
        margin: .only(left: 7.px),
        color: textMutedColor,
      ),
      css('.pagination').styles(
        display: .flex,
        margin: .only(top: 54.px),
        justifyContent: .center,
        alignItems: .center,
        gap: .all(8.px),
      ),
      css('.pagination button', [
        css('&').styles(
          display: .flex,
          width: 34.px,
          height: 34.px,
          padding: .zero,
          border: .all(style: .solid, color: borderColor, width: 1.px),
          radius: .all(.circular(7.px)),
          cursor: .pointer,
          justifyContent: .center,
          alignItems: .center,
          color: textSecondaryColor,
          fontFamily: .inherit,
          fontSize: 0.84.rem,
          backgroundColor: surfaceColor,
        ),
        css('&:hover:not(:disabled)').styles(
          border: .all(color: primaryColor, width: 1.px),
          color: primaryColor,
        ),
        css('&.active').styles(
          border: .all(color: primaryColor, width: 1.px),
          color: onPrimaryColor,
          backgroundColor: primaryColor,
        ),
        css('&:disabled').styles(opacity: 0.35, cursor: .defaultCursor),
        css('&:focus-visible').styles(
          outline: Outline(color: primaryColor, style: .solid, width: OutlineWidth(2.px), offset: 2.px),
        ),
      ]),
      css('.pagination .pagination-arrow').styles(fontSize: 1.rem),
    ]),
    css.media(MediaQuery.screen(maxWidth: 700.px), [
      css('.posts-page').styles(
        padding: .only(top: 112.px, right: 28.px, bottom: 52.px, left: 28.px),
      ),
      css('.posts-page .posts-container').styles(maxWidth: 100.percent),
      css('.posts-page .posts-avatar').styles(width: 98.px, height: 98.px),
      css('.posts-page .posts-author-name').styles(fontSize: 2.25.rem),
      css('.posts-page .posts-tagline').styles(fontSize: 0.91.rem),
      css('.posts-page .archive-list').styles(gap: .all(38.px)),
      css('.posts-page .archive-post h2').styles(fontSize: 1.04.rem),
    ]),
  ];
}

class _PostsPageState extends State<PostsPage> {
  int _currentPage = 1;

  List<PostPreview> get _posts => component.posts.map(PostPreview.fromMap).toList(growable: false);

  void _goToPage(int page) {
    final nextPage = page.clamp(1, (_posts.length / component.pageSize).ceil());
    if (nextPage == _currentPage) return;
    setState(() => _currentPage = nextPage);
  }

  @override
  Component build(BuildContext context) {
    final result = paginateItems(_posts, page: _currentPage, pageSize: component.pageSize);
    final hasProfile =
        component.avatarUrl.isNotEmpty || component.authorName.isNotEmpty || component.tagline.isNotEmpty;

    return section(classes: 'posts-page', [
      main_(classes: 'posts-container', [
        if (hasProfile)
          header(classes: 'posts-profile', [
            if (component.avatarUrl.isNotEmpty)
              img(
                src: component.avatarUrl,
                alt: component.authorName.isEmpty ? 'Ảnh đại diện' : 'Ảnh đại diện của ${component.authorName}',
                width: 100,
                height: 100,
                classes: 'posts-avatar',
              ),
            if (component.authorName.isNotEmpty) h1(classes: 'posts-author-name', [.text(component.authorName)]),
            if (component.tagline.isNotEmpty) p(classes: 'posts-tagline', [.text(component.tagline)]),
          ]),
        if (_posts.isEmpty)
          p(classes: 'empty-posts-message', [.text('Không có bài viết nào.')])
        else ...[
          div(classes: 'archive-list', [
            for (final post in result.items) _ArchivePost(post: post),
          ]),
          nav(
            classes: 'pagination',
            attributes: const {'aria-label': 'Phân trang bài viết'},
            [
              button(
                type: .button,
                classes: 'pagination-arrow',
                disabled: !result.hasPreviousPage,
                attributes: const {'aria-label': 'Trang trước'},
                onClick: () => _goToPage(result.currentPage - 1),
                [.text('‹')],
              ),
              for (var page = 1; page <= result.totalPages; page++)
                button(
                  type: .button,
                  classes: page == result.currentPage ? 'active' : null,
                  attributes: {
                    'aria-label': 'Trang $page',
                    if (page == result.currentPage) 'aria-current': 'page',
                  },
                  onClick: () => _goToPage(page),
                  [.text('$page')],
                ),
              button(
                type: .button,
                classes: 'pagination-arrow',
                disabled: !result.hasNextPage,
                attributes: const {'aria-label': 'Trang sau'},
                onClick: () => _goToPage(result.currentPage + 1),
                [.text('›')],
              ),
            ],
          ),
        ],
      ]),
    ]);
  }
}

class _ArchivePost extends StatelessComponent {
  const _ArchivePost({required this.post});

  final PostPreview post;

  @override
  Component build(BuildContext context) {
    final Component title = post.href == null ? .text(post.title) : a(href: post.href!, [.text(post.title)]);

    return article(classes: 'archive-post post-hover-card', [
      h2([title]),
      p(classes: 'archive-excerpt', [.text(post.excerpt)]),
      div(classes: 'archive-meta', [
        span(classes: 'archive-meta-item', [.text(post.publishedAt)]),
        span(classes: 'archive-meta-item', [.text(post.category)]),
        for (final tag in post.tags)
          span(classes: 'archive-meta-item', [
            a(href: tag.href, [.text('#${tag.name}')]),
          ]),
      ]),
    ]);
  }
}
