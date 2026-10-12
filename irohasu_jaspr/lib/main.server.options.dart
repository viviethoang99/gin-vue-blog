// dart format off
// ignore_for_file: type=lint

// GENERATED FILE, DO NOT MODIFY
// Generated with jaspr_builder

import 'package:jaspr/server.dart';
import 'package:irohasu_jaspr/components/article_search.dart'
    as _article_search;
import 'package:irohasu_jaspr/components/article_search_trigger.dart'
    as _article_search_trigger;
import 'package:irohasu_jaspr/components/header.dart' as _header;
import 'package:irohasu_jaspr/components/site_footer.dart' as _site_footer;
import 'package:irohasu_jaspr/components/static_particle_background.dart'
    as _static_particle_background;
import 'package:irohasu_jaspr/components/theme_toggle.dart' as _theme_toggle;
import 'package:irohasu_jaspr/components/toast.dart' as _toast;
import 'package:irohasu_jaspr/constants/theme.dart' as _theme;
import 'package:irohasu_jaspr/pages/about.dart' as _about;
import 'package:irohasu_jaspr/pages/home.dart' as _home;
import 'package:irohasu_jaspr/pages/post_detail.dart' as _post_detail;
import 'package:irohasu_jaspr/pages/posts.dart' as _posts;
import 'package:irohasu_jaspr/pages/tag_articles.dart' as _tag_articles;
import 'package:irohasu_jaspr/pages/tags.dart' as _tags;
import 'package:irohasu_jaspr/app.dart' as _app;

/// Default [ServerOptions] for use with your Jaspr project.
///
/// Use this to initialize Jaspr **before** calling [runApp].
///
/// Example:
/// ```dart
/// import 'main.server.options.dart';
///
/// void main() {
///   Jaspr.initializeApp(
///     options: defaultServerOptions,
///   );
///
///   runApp(...);
/// }
/// ```
ServerOptions get defaultServerOptions => ServerOptions(
  clientId: 'main.client.dart.js',
  clients: {
    _article_search.ArticleSearchOverlay:
        ClientTarget<_article_search.ArticleSearchOverlay>(
          'article_search',
          params: __article_searchArticleSearchOverlay,
        ),
    _article_search_trigger.ArticleSearchTrigger:
        ClientTarget<_article_search_trigger.ArticleSearchTrigger>(
          'article_search_trigger',
          params: __article_search_triggerArticleSearchTrigger,
        ),
    _static_particle_background.StaticParticleBackground:
        ClientTarget<_static_particle_background.StaticParticleBackground>(
          'static_particle_background',
          params: __static_particle_backgroundStaticParticleBackground,
        ),
    _theme_toggle.ThemeToggle: ClientTarget<_theme_toggle.ThemeToggle>(
      'theme_toggle',
    ),
    _toast.ToastViewport: ClientTarget<_toast.ToastViewport>(
      'toast',
      params: __toastToastViewport,
    ),
    _about.About: ClientTarget<_about.About>('about'),
    _post_detail.TableOfContents: ClientTarget<_post_detail.TableOfContents>(
      'post_detail',
      params: __post_detailTableOfContents,
    ),
    _posts.PostsPage: ClientTarget<_posts.PostsPage>(
      'posts',
      params: __postsPostsPage,
    ),
    _tags.TagsPage: ClientTarget<_tags.TagsPage>(
      'tags',
      params: __tagsTagsPage,
    ),
  },
  styles: () => [
    ..._theme.styles,
    ..._app.App.styles,
    ..._article_search.ArticleSearchOverlay.styles,
    ..._header.Header.styles,
    ..._site_footer.SiteFooter.styles,
    ..._static_particle_background.StaticParticleBackground.styles,
    ..._theme_toggle.ThemeToggle.styles,
    ..._toast.ToastViewport.styles,
    ..._about.About.styles,
    ..._home.Home.styles,
    ..._post_detail.PostDetailPage.styles,
    ..._posts.PostsPage.styles,
    ..._tag_articles.TagArticlesPage.styles,
    ..._tags.TagsPage.styles,
  ],
);

Map<String, Object?> __article_searchArticleSearchOverlay(
  _article_search.ArticleSearchOverlay c,
) => {'apiBaseUrl': c.apiBaseUrl};
Map<String, Object?> __article_search_triggerArticleSearchTrigger(
  _article_search_trigger.ArticleSearchTrigger c,
) => {'classes': c.classes};
Map<String, Object?> __static_particle_backgroundStaticParticleBackground(
  _static_particle_background.StaticParticleBackground c,
) => {
  'nodeCount': c.nodeCount,
  'maxNodeCount': c.maxNodeCount,
  'nodeCountReferenceArea': c.nodeCountReferenceArea,
  'connectionProbability': c.connectionProbability,
  'connectionDistanceScale': c.connectionDistanceScale,
  'minConnectionDistance': c.minConnectionDistance,
  'maxConnectionDistance': c.maxConnectionDistance,
  'viewportOverscan': c.viewportOverscan,
  'lineWidth': c.lineWidth,
  'lightColor': c.lightColor,
  'darkColor': c.darkColor,
  'lightMinOpacity': c.lightMinOpacity,
  'lightMaxOpacity': c.lightMaxOpacity,
  'darkMinOpacity': c.darkMinOpacity,
  'darkMaxOpacity': c.darkMaxOpacity,
  'animationDurationMs': c.animationDurationMs,
  'resizeDebounceMs': c.resizeDebounceMs,
};
Map<String, Object?> __toastToastViewport(_toast.ToastViewport c) => {
  'maxVisible': c.maxVisible,
};
Map<String, Object?> __post_detailTableOfContents(
  _post_detail.TableOfContents c,
) => {'headings': c.headings};
Map<String, Object?> __postsPostsPage(_posts.PostsPage c) => {
  'posts': c.posts,
  'authorName': c.authorName,
  'tagline': c.tagline,
  'avatarUrl': c.avatarUrl,
  'pageSize': c.pageSize,
};
Map<String, Object?> __tagsTagsPage(_tags.TagsPage c) => {'tags': c.tags};
