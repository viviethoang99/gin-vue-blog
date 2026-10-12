import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';

import '../constants/theme.dart';
import '../services/article_metrics.dart';

class TagArticleData {
  const TagArticleData({
    required this.id,
    required this.title,
    required this.excerpt,
    required this.publishedAt,
    required this.category,
    required this.wordCount,
    this.imageUrl,
  });

  final int id;
  final String title;
  final String excerpt;
  final String publishedAt;
  final String category;
  final int wordCount;
  final String? imageUrl;
}

class TagArticlesPage extends StatelessComponent {
  const TagArticlesPage({required this.tagName, required this.posts, super.key});

  final String tagName;
  final List<TagArticleData> posts;

  @override
  Component build(BuildContext context) {
    return section(classes: 'tag-articles-page', [
      main_(classes: 'tag-articles-container', [
        header(classes: 'tag-articles-header', [
          p(classes: 'tag-articles-eyebrow', [.text('Thẻ')]),
          h1([.text(tagName)]),
          p(classes: 'tag-articles-summary', [.text('Tổng cộng ${posts.length} bài viết')]),
        ]),
        if (posts.isEmpty)
          p(classes: 'tag-articles-empty', [.text('Chưa có bài viết nào thuộc tag này.')])
        else
          div(classes: 'tag-article-list', [
            for (final post in posts) _TagArticleCard(post: post),
          ]),
      ]),
    ]);
  }

  @css
  static List<StyleRule> get styles => [
    css('.tag-articles-page').styles(
      width: 100.percent,
      minHeight: 100.vh,
      padding: .only(top: 92.px, right: 80.px, bottom: 70.px, left: 80.px),
      margin: .zero,
      boxSizing: .borderBox,
      color: textSecondaryColor,
      backgroundColor: Colors.transparent,
    ),
    css('.tag-articles-container').styles(
      width: 100.percent,
      maxWidth: 960.px,
      margin: .only(left: .auto, right: .auto),
    ),
    css('.tag-articles-header').styles(margin: .only(bottom: 18.px)),
    css('.tag-articles-eyebrow').styles(
      color: textMutedColor,
      fontSize: 0.8.rem,
      lineHeight: 1.4.em,
    ),
    css('.tag-articles-header h1').styles(
      margin: .only(top: 5.px),
      color: textPrimaryColor,
      fontSize: 1.85.rem,
      fontWeight: .w600,
      lineHeight: 1.25.em,
    ),
    css('.tag-articles-summary').styles(
      margin: .only(top: 7.px),
      color: textMutedColor,
      fontSize: 0.88.rem,
      lineHeight: 1.5.em,
    ),
    css('.tag-article-card').styles(
      display: .flex,
      padding: .symmetric(vertical: 26.px),
      border: .only(
        bottom: .solid(color: borderColor, width: 1.px),
      ),
      alignItems: .center,
      gap: .all(40.px),
    ),
    css('.tag-article-content').styles(minWidth: 0.px, flex: Flex(grow: 1)),
    css('.tag-article-meta').styles(
      display: .flex,
      margin: .only(bottom: 9.px),
      flexWrap: .wrap,
      gap: .all(7.px),
      color: textMutedColor,
      fontSize: 0.78.rem,
      lineHeight: 1.45.em,
    ),
    css('.tag-article-meta-item:not(:last-child)::after').styles(
      content: '·',
      margin: .only(left: 7.px),
      color: textMutedColor,
    ),
    css('.tag-article-title').styles(
      margin: .zero,
      color: textPrimaryColor,
      fontSize: 1.18.rem,
      fontWeight: .w600,
      lineHeight: 1.45.em,
    ),
    css('.tag-article-title a').styles(color: .inherit, textDecoration: .none),
    css('.tag-article-title a:hover, .tag-article-title a:focus-visible').styles(color: primaryColor),
    css('.tag-article-excerpt').styles(
      margin: .only(top: 9.px),
      color: textSecondaryColor,
      fontSize: 0.9.rem,
      lineHeight: 1.72.em,
    ),
    css('.tag-article-image-link').styles(
      display: .block,
      width: 300.px,
      height: 168.px,
      radius: .all(.circular(7.px)),
      overflow: .hidden,
      flex: .none,
      backgroundColor: surfaceColor,
    ),
    css('.tag-article-image').styles(
      display: .block,
      width: 100.percent,
      height: 100.percent,
      raw: const {'object-fit': 'cover'},
    ),
    css('.tag-articles-empty').styles(
      padding: .only(top: 28.px),
      color: textMutedColor,
      fontSize: 0.92.rem,
    ),
    css.media(MediaQuery.screen(maxWidth: 900.px), [
      css('.tag-articles-page').styles(
        padding: .only(top: 104.px, right: 44.px, bottom: 60.px, left: 96.px),
      ),
      css('.tag-articles-page .tag-article-card').styles(gap: .all(28.px)),
      css('.tag-articles-page .tag-article-image-link').styles(width: 240.px, height: 144.px),
    ]),
    css.media(MediaQuery.screen(maxWidth: 700.px), [
      css('.tag-articles-page').styles(
        padding: .only(top: 104.px, right: 28.px, bottom: 52.px, left: 28.px),
      ),
      css('.tag-articles-page .tag-articles-header h1').styles(fontSize: 1.65.rem),
      css('.tag-articles-page .tag-article-card').styles(
        padding: .symmetric(vertical: 24.px),
        flexDirection: .column,
        alignItems: .stretch,
        gap: .all(18.px),
      ),
      css('.tag-articles-page .tag-article-image-link').styles(
        width: 100.percent,
        height: .auto,
        aspectRatio: const AspectRatio(16, 9),
        order: -1,
      ),
      css('.tag-articles-page .tag-article-title').styles(fontSize: 1.08.rem),
      css('.tag-articles-page .tag-article-excerpt').styles(fontSize: 0.88.rem, lineHeight: 1.65.em),
    ]),
  ];
}

class _TagArticleCard extends StatelessComponent {
  const _TagArticleCard({required this.post});

  final TagArticleData post;

  @override
  Component build(BuildContext context) {
    final href = '/posts/${post.id}';
    return article(classes: 'tag-article-card', [
      div(classes: 'tag-article-content', [
        div(classes: 'tag-article-meta', [
          span(classes: 'tag-article-meta-item', [.text(post.category)]),
          span(classes: 'tag-article-meta-item', [.text(post.publishedAt)]),
          span(classes: 'tag-article-meta-item', [
            .text('Khoảng ${formatVietnameseInteger(post.wordCount)} từ'),
          ]),
        ]),
        h2(classes: 'tag-article-title', [
          a(href: href, [.text(post.title)]),
        ]),
        if (post.excerpt.isNotEmpty) p(classes: 'tag-article-excerpt', [.text(post.excerpt)]),
      ]),
      if (post.imageUrl case final imageUrl?)
        a(href: href, classes: 'tag-article-image-link', [
          img(
            src: imageUrl,
            alt: 'Ảnh minh họa cho ${post.title}',
            width: 300,
            height: 168,
            loading: .lazy,
            classes: 'tag-article-image',
          ),
        ]),
    ]);
  }
}
