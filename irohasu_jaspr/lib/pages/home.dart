import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';

import '../constants/theme.dart';
import '../models/home_page_data.dart';

class Home extends StatelessComponent {
  const Home({required this.data, super.key});

  final HomePageData data;

  @override
  Component build(BuildContext context) {
    return section(classes: 'home-page', [
      main_(classes: 'home-layout', [
        if (_hasProfileData(data.profile)) _ProfileIntroduction(profile: data.profile),
        _RecentPosts(posts: data.recentPosts),
      ]),
    ]);
  }

  @css
  static List<StyleRule> get styles => [
    css('.home-page', [
      css('&').styles(
        display: .flex,
        width: 100.percent,
        minHeight: 100.vh,
        padding: .only(top: 112.px, right: 80.px, bottom: 24.px, left: 80.px),
        margin: .zero,
        boxSizing: .borderBox,
        flexDirection: .column,
        alignItems: .stretch,
        color: textSecondaryColor,
        backgroundColor: Colors.transparent,
      ),
      css('a', [
        css('&').styles(
          color: primaryColor,
          textDecoration: TextDecoration(line: .none),
        ),
        css('&:hover').styles(textDecoration: TextDecoration(line: .underline)),
        css('&:focus-visible').styles(
          radius: .all(.circular(3.px)),
          outline: Outline(color: primaryColor, style: .solid, width: OutlineWidth(2.px), offset: 3.px),
        ),
      ]),
      css('.home-layout').styles(
        display: .flex,
        width: 100.percent,
        maxWidth: 1220.px,
        margin: .only(left: .auto, right: .auto),
        flexDirection: .row,
        alignItems: .start,
        gap: .all(70.px),
        flex: const Flex(grow: 1.35, shrink: 1),
      ),
      css('.profile-column').styles(
        minWidth: 0.px,
        flex: const Flex(grow: 1.35, shrink: 1, basis: Unit.zero),
      ),
      css('.profile-avatar').styles(
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
      css('.profile-name').styles(
        margin: .only(top: 12.px),
        color: textPrimaryColor,
        fontFamily: const .list([FontFamily('Caveat'), FontFamilies.cursive]),
        fontSize: 2.4.rem,
        fontWeight: .w600,
        letterSpacing: 1.px,
        lineHeight: 1.em,
      ),
      css('.profile-tagline').styles(
        margin: .only(top: 7.px, bottom: 20.px),
        color: textMutedColor,
        fontSize: 0.96.rem,
        lineHeight: 1.5.em,
      ),
      css('.profile-biography p').styles(
        margin: .only(bottom: 14.px),
        color: textSecondaryColor,
        fontSize: 0.98.rem,
        lineHeight: 1.62.em,
      ),
      css('.posts-column').styles(
        minWidth: 0.px,
        flex: const Flex(grow: 0.85, shrink: 1, basis: Unit.zero),
      ),
      css('.section-heading').styles(
        position: .relative(),
        padding: .only(bottom: 8.px),
        margin: .only(bottom: 28.px),
        border: .only(
          bottom: BorderSide.dashed(color: borderColor, width: 1.px),
        ),
        color: textPrimaryColor,
        fontSize: 1.18.rem,
        fontWeight: .w600,
      ),
      css('.section-heading::after').styles(
        content: '',
        display: .block,
        position: .absolute(bottom: (-1).px, left: 0.px),
        width: 34.px,
        height: 2.px,
        backgroundColor: primaryColor,
      ),
      css('.post-list').styles(
        display: .flex,
        flexDirection: .column,
        gap: .all(34.px),
      ),
      css('.empty-posts-message').styles(
        margin: .zero,
        color: textMutedColor,
        fontSize: 0.9.rem,
        lineHeight: 1.55.em,
      ),
      css('.post-preview h3').styles(
        margin: .only(bottom: 6.px),
        color: textPrimaryColor,
        fontSize: 1.08.rem,
        fontWeight: .w600,
        lineHeight: 1.35.em,
      ),
      css('.post-preview h3 a').styles(color: textPrimaryColor),
      css('.post-excerpt').styles(
        margin: .zero,
        color: textSecondaryColor,
        fontSize: 0.9.rem,
        lineHeight: 1.55.em,
        raw: const {
          'overflow-wrap': 'anywhere',
          'word-break': 'break-word',
        },
      ),
      css('.post-meta').styles(
        display: .flex,
        margin: .only(top: 6.px),
        flexWrap: .wrap,
        gap: Gap(row: 3.px, column: 7.px),
        color: primaryColor,
        fontSize: 0.9.rem,
        lineHeight: 1.4.em,
      ),
      css('.post-meta-item:not(:last-child)::after').styles(
        content: '·',
        margin: .only(left: 7.px),
        color: textMutedColor,
      ),
      css('.more-posts').styles(
        display: .inlineBlock,
        margin: .only(top: 30.px),
        fontSize: 0.88.rem,
        fontWeight: .w600,
      ),
    ]),
    css.media(MediaQuery.screen(maxWidth: 1360.px), [
      css('.home-page').styles(
        padding: .only(top: 112.px, right: 96.px, bottom: 24.px, left: 126.px),
      ),
      css('.home-page .home-layout').styles(
        maxWidth: 740.px,
        flexDirection: .column,
        gap: .all(30.px),
      ),
      css('.home-page .profile-column, .home-page .posts-column').styles(
        width: 100.percent,
        flex: .none,
      ),
      css('.home-page .section-heading').styles(
        margin: .only(top: 4.px, bottom: 28.px),
      ),
    ]),
    css.media(MediaQuery.screen(maxWidth: 700.px), [
      css('.home-page').styles(
        padding: .only(top: 112.px, right: 28.px, bottom: 40.px, left: 28.px),
      ),
      css('.home-page .home-layout').styles(maxWidth: 100.percent),
      css('.home-page .profile-avatar').styles(width: 98.px, height: 98.px),
      css('.home-page .profile-name').styles(fontSize: 2.25.rem),
      css('.home-page .profile-tagline').styles(fontSize: 0.91.rem),
      css('.home-page .profile-biography p').styles(fontSize: 0.93.rem, lineHeight: 1.58.em),
      css('.home-page .post-list').styles(gap: .all(38.px)),
      css('.home-page .post-preview h3').styles(fontSize: 1.04.rem),
    ]),
  ];
}

class _ProfileIntroduction extends StatelessComponent {
  const _ProfileIntroduction({required this.profile});

  final AuthorProfile profile;

  @override
  Component build(BuildContext context) {
    return div(classes: 'profile-column', [
      if (profile.avatarUrl.isNotEmpty)
        img(
          src: profile.avatarUrl,
          alt: profile.name.isEmpty ? 'Ảnh đại diện' : 'Ảnh đại diện của ${profile.name}',
          width: 100,
          height: 100,
          classes: 'profile-avatar',
        ),
      if (profile.name.isNotEmpty) h1(classes: 'profile-name', [.text(profile.name)]),
      if (profile.tagline.isNotEmpty) p(classes: 'profile-tagline', [.text(profile.tagline)]),
      if (profile.biography.isNotEmpty)
        div(classes: 'profile-biography', [
          for (final paragraph in profile.biography) p(_renderSegments(paragraph.segments)),
        ]),
    ]);
  }
}

class _RecentPosts extends StatelessComponent {
  const _RecentPosts({required this.posts});

  final List<PostPreview> posts;

  @override
  Component build(BuildContext context) {
    return section(id: 'recent-posts', classes: 'posts-column', [
      h2(classes: 'section-heading', [.text('Bài viết gần đây')]),
      if (posts.isEmpty)
        p(classes: 'empty-posts-message', [.text('Không có bài viết nào.')])
      else ...[
        div(classes: 'post-list', [
          for (final post in posts) _PostPreviewCard(post: post),
        ]),
        a(href: '/posts', classes: 'more-posts', [.text('» Xem thêm «')]),
      ],
    ]);
  }
}

class _PostPreviewCard extends StatelessComponent {
  const _PostPreviewCard({required this.post});

  final PostPreview post;

  @override
  Component build(BuildContext context) {
    final Component title = post.href == null ? .text(post.title) : a(href: post.href!, [.text(post.title)]);

    return article(classes: 'post-preview post-hover-card', [
      h3([title]),
      p(classes: 'post-excerpt', [.text(post.excerpt)]),
      div(classes: 'post-meta', [
        span(classes: 'post-meta-item', [.text(post.publishedAt)]),
        span(classes: 'post-meta-item', [.text(post.category)]),
        for (final tag in post.tags)
          span(classes: 'post-meta-item', [
            a(href: tag.href, [.text('#${tag.name}')]),
          ]),
      ]),
    ]);
  }
}

bool _hasProfileData(AuthorProfile profile) =>
    profile.name.isNotEmpty ||
    profile.tagline.isNotEmpty ||
    profile.avatarUrl.isNotEmpty ||
    profile.biography.isNotEmpty;

List<Component> _renderSegments(List<TextSegment> segments) => [
  for (final segment in segments)
    if (segment.href case final href?) a(href: href, [.text(segment.text)]) else .text(segment.text),
];
