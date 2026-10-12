// dart format off
// ignore_for_file: type=lint

// GENERATED FILE, DO NOT MODIFY
// Generated with jaspr_builder

import 'package:jaspr/client.dart';

import 'package:irohasu_jaspr/components/article_search.dart'
    deferred as _article_search;
import 'package:irohasu_jaspr/components/article_search_trigger.dart'
    deferred as _article_search_trigger;
import 'package:irohasu_jaspr/components/static_particle_background.dart'
    deferred as _static_particle_background;
import 'package:irohasu_jaspr/components/theme_toggle.dart'
    deferred as _theme_toggle;
import 'package:irohasu_jaspr/components/toast.dart' deferred as _toast;
import 'package:irohasu_jaspr/pages/about.dart' deferred as _about;
import 'package:irohasu_jaspr/pages/post_detail.dart' deferred as _post_detail;
import 'package:irohasu_jaspr/pages/posts.dart' deferred as _posts;
import 'package:irohasu_jaspr/pages/tags.dart' deferred as _tags;

/// Default [ClientOptions] for use with your Jaspr project.
///
/// Use this to initialize Jaspr **before** calling [runApp].
///
/// Example:
/// ```dart
/// import 'main.client.options.dart';
///
/// void main() {
///   Jaspr.initializeApp(
///     options: defaultClientOptions,
///   );
///
///   runApp(...);
/// }
/// ```
ClientOptions get defaultClientOptions => ClientOptions(
  clients: {
    'article_search': ClientLoader(
      (p) => _article_search.ArticleSearchOverlay(
        apiBaseUrl: p['apiBaseUrl'] as String,
      ),
      loader: _article_search.loadLibrary,
    ),
    'article_search_trigger': ClientLoader(
      (p) => _article_search_trigger.ArticleSearchTrigger(
        classes: p['classes'] as String?,
      ),
      loader: _article_search_trigger.loadLibrary,
    ),
    'static_particle_background': ClientLoader(
      (p) => _static_particle_background.StaticParticleBackground(
        nodeCount: p['nodeCount'] as int,
        maxNodeCount: p['maxNodeCount'] as int,
        nodeCountReferenceArea: p['nodeCountReferenceArea'] as double,
        connectionProbability: p['connectionProbability'] as double,
        connectionDistanceScale: p['connectionDistanceScale'] as double,
        minConnectionDistance: p['minConnectionDistance'] as double,
        maxConnectionDistance: p['maxConnectionDistance'] as double,
        viewportOverscan: p['viewportOverscan'] as double,
        lineWidth: p['lineWidth'] as double,
        lightColor: p['lightColor'] as String,
        darkColor: p['darkColor'] as String,
        lightMinOpacity: p['lightMinOpacity'] as double,
        lightMaxOpacity: p['lightMaxOpacity'] as double,
        darkMinOpacity: p['darkMinOpacity'] as double,
        darkMaxOpacity: p['darkMaxOpacity'] as double,
        animationDurationMs: p['animationDurationMs'] as int,
        resizeDebounceMs: p['resizeDebounceMs'] as int,
      ),
      loader: _static_particle_background.loadLibrary,
    ),
    'theme_toggle': ClientLoader(
      (p) => _theme_toggle.ThemeToggle(),
      loader: _theme_toggle.loadLibrary,
    ),
    'toast': ClientLoader(
      (p) => _toast.ToastViewport(maxVisible: p['maxVisible'] as int),
      loader: _toast.loadLibrary,
    ),
    'about': ClientLoader((p) => _about.About(), loader: _about.loadLibrary),
    'post_detail': ClientLoader(
      (p) => _post_detail.TableOfContents(
        headings: (p['headings'] as List<Object?>)
            .map((i) => (i as Map<String, Object?>))
            .toList(),
      ),
      loader: _post_detail.loadLibrary,
    ),
    'posts': ClientLoader(
      (p) => _posts.PostsPage(
        posts: (p['posts'] as List<Object?>)
            .map((i) => (i as Map<String, Object?>))
            .toList(),
        authorName: p['authorName'] as String,
        tagline: p['tagline'] as String,
        avatarUrl: p['avatarUrl'] as String,
        pageSize: p['pageSize'] as int,
      ),
      loader: _posts.loadLibrary,
    ),
    'tags': ClientLoader(
      (p) => _tags.TagsPage(
        tags: (p['tags'] as List<Object?>)
            .map((i) => (i as Map<String, Object?>))
            .toList(),
      ),
      loader: _tags.loadLibrary,
    ),
  },
);
