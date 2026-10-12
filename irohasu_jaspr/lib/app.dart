import 'package:jaspr/dom.dart';
import 'package:jaspr_router/jaspr_router.dart';
import 'package:jaspr/server.dart';

import 'components/header.dart';
import 'components/site_footer.dart';
import 'components/static_particle_background.dart';
import 'components/toast.dart';
import 'config/environment.dart';
import 'data/model/post.dart';
import 'data/post_api/post_mapper.dart';
import 'data/post_api/post_repository.dart';
import 'data/post_api/home_data.dart';
import 'data/post_api/home_mapper.dart';
import 'data/site_navigation_data.dart';
import 'models/home_page_data.dart';
import 'pages/about.dart';
import 'pages/home.dart';
import 'pages/post_detail.dart';
import 'pages/posts.dart';
import 'pages/tag_articles.dart';
import 'pages/tags.dart';
import 'services/article_metrics.dart';
import 'services/tag_summary.dart';

// The server-rendered application shell. Interactive @client components are
// hydrated in the browser after each route is rendered.
class App extends AsyncStatelessComponent {
  const App({super.key});

  @override
  Future<Component> build(BuildContext context) async {
    final postRepository = PostRepository.live();
    final (posts, remoteHome) = await (
      _loadPosts(postRepository),
      _loadHome(postRepository),
    ).wait;
    final previews = [
      for (final post in posts) post.toPreview(href: post.id == null ? null : '/posts/${post.id}'),
    ];
    final recentPosts = previews.take(3).toList(growable: false);
    final homeData = remoteHome?.toPageData(recentPosts: recentPosts) ?? HomePageData.withRecentPosts(recentPosts);
    final tags = collectTagSummaries(posts);
    final tagIdsByName = {
      for (final tag in tags) tag.name.trim().toLowerCase(): tag.id,
    };
    final tagRoutes = <int, String>{
      for (final post in posts)
        for (final tag in post.tags ?? const <Category>[])
          if (tag.name.trim().isNotEmpty) tag.id: tag.name.trim(),
    };

    return .fragment([
      const StaticParticleBackground(),
      const ToastViewport(),
      div(classes: 'main', [
        Header(
          avatarUrl: homeData.profile.avatarUrl,
          items: siteNavigationItems,
        ),
        div(classes: 'route-content', [
          Router(
            routes: [
              Route(
                path: '/',
                title: 'Irohasu',
                builder: (context, state) => Home(data: homeData),
              ),
              Route(
                path: '/posts',
                title: 'Bài viết | Irohasu',
                builder: (context, state) => PostsPage(
                  posts: [for (final preview in previews) preview.toMap()],
                  authorName: homeData.profile.name,
                  tagline: homeData.profile.tagline,
                  avatarUrl: homeData.profile.avatarUrl,
                ),
              ),
              Route(
                path: '/tags',
                title: 'Tag | Irohasu',
                builder: (context, state) => TagsPage(
                  tags: [
                    for (final tag in tags) {'id': tag.id, 'name': tag.name, 'count': tag.count},
                  ],
                ),
              ),
              for (final tag in tagRoutes.entries)
                Route(
                  path: '/tags/${tag.key}',
                  title: '${tag.value} | Irohasu',
                  builder: (context, state) => TagArticlesPage(
                    tagName: tag.value,
                    posts: [
                      for (final post in posts)
                        if (post.id != null && _hasTag(post, tag.value)) _toTagArticle(post),
                    ],
                  ),
                ),
              Route(
                path: '/posts/:id',
                title: 'Bài viết | Irohasu',
                builder: (context, state) {
                  final id = int.tryParse(state.params['id'] ?? '');
                  final fallback = posts.where((post) => post.id == id).firstOrNull;
                  return _PostDetailLoader(
                    repository: postRepository,
                    id: id,
                    fallback: fallback,
                    tagIdsByName: tagIdsByName,
                  );
                },
              ),
              Route(path: '/about', title: 'About', builder: (context, state) => const About()),
            ],
          ),
        ]),
        SiteFooter(data: homeData.footer),
      ]),
    ]);
  }

  // Defines the CSS styles for this component.
  //
  // By using the @css annotation, these will be rendered automatically to CSS and included in your page.
  // Must be a variable or getter of type [List<StyleRule>].
  @css
  static List<StyleRule> get styles => [
    css('.main', [
      // The '&' refers to the parent selector of a nested style rules.
      css('&').styles(
        display: .flex,
        position: .relative(),
        zIndex: ZIndex(1),
        minHeight: 100.vh,
        flexDirection: .column,
      ),
      css('.route-content').styles(
        display: .flex,
        width: 100.percent,
        minHeight: 0.px,
        flex: .grow(1),
      ),
      css(
        '.route-content > .home-page, .route-content > .posts-page, .route-content > .post-detail-page, .route-content > .tags-page, .route-content > .tag-articles-page, .route-content > .about-page',
      ).styles(
        minHeight: 0.px,
        flex: .grow(1),
      ),
    ]),
  ];
}

class _PostDetailLoader extends AsyncStatelessComponent {
  const _PostDetailLoader({
    required this.repository,
    required this.id,
    required this.fallback,
    required this.tagIdsByName,
  });

  final PostRepository repository;
  final int? id;
  final Post? fallback;
  final Map<String, int> tagIdsByName;

  @override
  Future<Component> build(BuildContext context) async {
    var post = fallback;
    final postId = id;
    if (postId != null) {
      try {
        final response = await repository.getPost(postId);
        if (response.code == 0) post = response.data;
      } catch (_) {
        // The list response remains a fallback if the detail request fails.
      }
    }
    if (post == null) {
      return const _PostNotFound();
    }
    return PostDetailPage(
      post: post.toPostDetail(
        slug: '${post.id ?? ''}',
        tagIdsByName: tagIdsByName,
      ),
    );
  }
}

Future<List<Post>> _loadPosts(PostRepository repository) async {
  try {
    const pageSize = 100;
    final posts = <Post>[];
    for (var page = 1; ; page++) {
      final response = await repository.getPosts(page: page, limit: pageSize);
      if (response.code != 0) break;
      final batch = response.data ?? const <Post>[];
      posts.addAll(batch);
      final total = response.paging?.total;
      if (batch.length < pageSize || (total != null && posts.length >= total)) break;
    }
    return posts;
  } catch (_) {
    // Render the empty state when the article API cannot be reached.
  }
  return const <Post>[];
}

Future<HomeData?> _loadHome(PostRepository repository) async {
  try {
    final response = await repository.getHome();
    if (response.code == 0) return response.data;
  } catch (_) {
    // Keep static generation available when the home API cannot be reached.
  }
  return null;
}

bool _hasTag(Post post, String tagName) {
  final normalizedName = tagName.trim().toLowerCase();
  return post.tags?.any((tag) => tag.name.trim().toLowerCase() == normalizedName) == true;
}

TagArticleData _toTagArticle(Post post) {
  final preview = post.toPreview(href: post.id == null ? null : '/posts/${post.id}');
  return TagArticleData(
    id: post.id!,
    title: preview.title,
    excerpt: preview.excerpt,
    publishedAt: preview.publishedAt,
    category: preview.category,
    wordCount: countArticleWords(post.content?.trim().isNotEmpty == true ? post.content! : post.desc ?? ''),
    imageUrl: _resolveArticleImage(post.img),
  );
}

String? _resolveArticleImage(String? value) {
  final source = value?.trim();
  if (source == null || source.isEmpty) return null;
  final uri = Uri.tryParse(source);
  if (uri == null) return null;
  if (uri.hasScheme) return uri.scheme == 'http' || uri.scheme == 'https' ? source : null;
  return Uri.parse(publicApiBaseUrl).resolveUri(uri).toString();
}

class _PostNotFound extends StatelessComponent {
  const _PostNotFound();

  @override
  Component build(BuildContext context) => main_(classes: 'post-detail-page', [
    h1([.text('Không tìm thấy bài viết')]),
    p([.text('Bài viết không tồn tại hoặc đã bị gỡ.')]),
  ]);
}
