import 'dart:async';

import 'package:dio/dio.dart';
import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';
import 'package:universal_web/web.dart' as web;

import '../config/environment.dart';
import '../constants/theme.dart';
import '../data/model/post.dart';
import '../data/post_api/post_repository.dart';

class ArticleSearchController {
  ArticleSearchController._();

  static final ArticleSearchController instance = ArticleSearchController._();

  final StreamController<void> _openEvents = StreamController<void>.broadcast();

  Stream<void> get openEvents => _openEvents.stream;

  void open() => _openEvents.add(null);
}

@client
class ArticleSearchOverlay extends StatefulComponent {
  const ArticleSearchOverlay({
    required this.apiBaseUrl,
    super.key,
  });

  final String apiBaseUrl;

  @override
  State<ArticleSearchOverlay> createState() => _ArticleSearchOverlayState();

  @css
  static List<StyleRule> get styles => [
    css('.article-search-overlay').styles(
      display: .flex,
      position: .fixed(top: 0.px, right: 0.px, bottom: 0.px, left: 0.px),
      zIndex: ZIndex(100),
      padding: .only(top: 12.vh, right: 24.px, bottom: 24.px, left: 24.px),
      boxSizing: .borderBox,
      overflow: .auto,
      backdropFilter: .blur(6.px),
      animation: Animation(name: 'article-search-backdrop-in', duration: 160.ms, curve: .easeOut),
      justifyContent: .center,
      alignItems: .start,
      backgroundColor: const Color('#12111899'),
    ),
    css('.article-search-dialog').styles(
      width: 620.px,
      maxWidth: 100.percent,
      maxHeight: 72.vh,
      padding: .all(18.px),
      boxSizing: .borderBox,
      border: .all(style: .solid, color: borderColor, width: 1.px),
      radius: .all(.circular(14.px)),
      overflow: .auto,
      shadow: BoxShadow(offsetX: 0.px, offsetY: 20.px, blur: 60.px, color: const Color('#0000004d')),
      animation: Animation(name: 'article-search-dialog-in', duration: 180.ms, curve: .easeOut),
      backgroundColor: const Color.variable('--color-surface-elevated'),
    ),
    css('.article-search-form').styles(
      display: .flex,
      alignItems: .center,
      gap: .all(10.px),
    ),
    css('.article-search-field').styles(
      minWidth: 0.px,
      padding: .symmetric(vertical: 12.px, horizontal: 14.px),
      border: .all(style: .solid, color: borderColor, width: 1.px),
      radius: .all(.circular(9.px)),
      flex: Flex(grow: 1),
      color: textPrimaryColor,
      fontFamily: .inherit,
      fontSize: 1.rem,
      backgroundColor: surfaceColor,
      raw: const {'outline': 'none'},
    ),
    css('.article-search-field:focus').styles(
      border: .all(style: .solid, color: primaryColor, width: 1.px),
      shadow: BoxShadow(offsetX: 0.px, offsetY: 0.px, blur: 0.px, spread: 2.px, color: pageAccentColor),
    ),
    css('.article-search-close').styles(
      display: .flex,
      width: 42.px,
      height: 42.px,
      padding: .zero,
      border: .none,
      radius: .all(.circular(8.px)),
      cursor: .pointer,
      justifyContent: .center,
      alignItems: .center,
      color: textSecondaryColor,
      fontSize: 1.35.rem,
      backgroundColor: Colors.transparent,
    ),
    css('.article-search-close:hover, .article-search-close:focus-visible').styles(
      color: primaryColor,
      backgroundColor: postHoverColor,
    ),
    css('.article-search-status').styles(
      padding: .only(top: 22.px, right: 4.px, bottom: 8.px, left: 4.px),
      color: textMutedColor,
      fontSize: 0.9.rem,
      lineHeight: 1.5.em,
    ),
    css('.article-search-status--error').styles(color: dangerColor),
    css('.article-search-results').styles(
      display: .flex,
      padding: .only(top: 14.px),
      flexDirection: .column,
      gap: .all(4.px),
    ),
    css('.article-search-result').styles(
      display: .block,
      padding: .symmetric(vertical: 12.px, horizontal: 14.px),
      radius: .all(.circular(9.px)),
      transition: Transition('background-color', duration: 150.ms, curve: .easeOut),
      color: textPrimaryColor,
      textDecoration: .none,
    ),
    css('.article-search-result:hover, .article-search-result:focus-visible').styles(
      backgroundColor: postHoverColor,
      raw: const {'outline': 'none'},
    ),
    css('.article-search-result-title').styles(
      color: textPrimaryColor,
      fontSize: 0.98.rem,
      fontWeight: .w600,
      lineHeight: 1.4.em,
    ),
    css('.article-search-result-excerpt').styles(
      margin: .only(top: 5.px),
      overflow: .hidden,
      color: textSecondaryColor,
      fontSize: 0.84.rem,
      lineHeight: 1.5.em,
      raw: const {
        'display': '-webkit-box',
        '-webkit-box-orient': 'vertical',
        '-webkit-line-clamp': '2',
      },
    ),
    css.keyframes('article-search-backdrop-in', {
      '0%': Styles(opacity: 0),
      '100%': Styles(opacity: 1),
    }),
    css.keyframes('article-search-dialog-in', {
      '0%': Styles(opacity: 0, transform: .translate(y: (-10).px)),
      '100%': Styles(opacity: 1, transform: .translate(y: 0.px)),
    }),
    css.media(MediaQuery.screen(maxWidth: 700.px), [
      css('.article-search-overlay').styles(
        padding: .only(top: 76.px, right: 16.px, bottom: 16.px, left: 16.px),
      ),
      css('.article-search-dialog').styles(maxHeight: .expression('calc(100vh - 92px)'), padding: .all(14.px)),
    ]),
    css.media(const MediaQuery.raw('(prefers-reduced-motion: reduce)'), [
      css('.article-search-overlay, .article-search-dialog').styles(raw: const {'animation': 'none'}),
    ]),
  ];
}

class _ArticleSearchOverlayState extends State<ArticleSearchOverlay> {
  final GlobalNodeKey<web.HTMLInputElement> _inputKey = GlobalNodeKey(debugLabel: 'article-search-input');
  StreamSubscription<void>? _openSubscription;
  StreamSubscription<web.Event>? _keySubscription;
  Timer? _debounce;
  late PostRepository _repository;
  List<Post> _results = const [];
  String _keyword = '';
  String? _error;
  bool _isOpen = false;
  bool _isLoading = false;
  int _requestId = 0;

  @override
  void initState() {
    super.initState();
    _repository = PostRepository.live(baseUrl: _normalizedBaseUrl(component.apiBaseUrl));
    _openSubscription = ArticleSearchController.instance.openEvents.listen((_) => _open());
    if (kIsWeb) {
      _keySubscription = web.EventStreamProviders.keyDownEvent.forTarget(web.window).listen((event) {
        if (_isOpen && event.key == 'Escape') _close();
      });
    }
  }

  void _open() {
    if (!mounted) return;
    setState(() => _isOpen = true);
    if (kIsWeb) web.document.body?.style.overflow = 'hidden';
    Future<void>.delayed(Duration.zero, () => _inputKey.currentNode?.focus());
  }

  void _close() {
    if (!_isOpen || !mounted) return;
    _debounce?.cancel();
    _requestId++;
    setState(() {
      _isOpen = false;
      _isLoading = false;
    });
    if (kIsWeb) web.document.body?.style.removeProperty('overflow');
  }

  void _onKeywordChanged(String value) {
    _keyword = value;
    _debounce?.cancel();
    if (value.trim().isEmpty) {
      _requestId++;
      setState(() {
        _results = const [];
        _error = null;
        _isLoading = false;
      });
      return;
    }
    _debounce = Timer(const Duration(milliseconds: 350), _search);
  }

  Future<void> _search() async {
    final keyword = _keyword.trim();
    if (keyword.isEmpty) return;
    final requestId = ++_requestId;
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final response = await _repository.searchPosts(keyword);
      if (!mounted || requestId != _requestId) return;
      setState(() {
        _isLoading = false;
        if (response.code == 0) {
          _results = response.data ?? const [];
        } else {
          _results = const [];
          _error = response.message.isEmpty ? 'Không thể tìm kiếm bài viết.' : response.message;
        }
      });
    } on DioException catch (_) {
      _showRequestError(requestId);
    } catch (_) {
      _showRequestError(requestId);
    }
  }

  void _showRequestError(int requestId) {
    if (!mounted || requestId != _requestId) return;
    setState(() {
      _isLoading = false;
      _results = const [];
      _error = 'Không thể kết nối đến máy chủ tìm kiếm. Vui lòng thử lại.';
    });
  }

  @override
  void dispose() {
    _debounce?.cancel();
    _openSubscription?.cancel();
    _keySubscription?.cancel();
    if (kIsWeb) web.document.body?.style.removeProperty('overflow');
    super.dispose();
  }

  @override
  Component build(BuildContext context) {
    if (!_isOpen) return .empty();

    return div(
      classes: 'article-search-overlay',
      attributes: const {'role': 'presentation'},
      events: {'click': (_) => _close()},
      [
        div(
          classes: 'article-search-dialog',
          attributes: const {
            'role': 'dialog',
            'aria-modal': 'true',
            'aria-label': 'Tìm kiếm bài viết',
          },
          events: {'click': (event) => event.stopPropagation()},
          [
            form(
              classes: 'article-search-form',
              events: {
                'submit': (event) {
                  event.preventDefault();
                  _debounce?.cancel();
                  _search();
                },
              },
              [
                input<String>(
                  key: _inputKey,
                  type: .search,
                  name: 'keyword',
                  value: _keyword,
                  classes: 'article-search-field',
                  attributes: const {
                    'placeholder': 'Tìm kiếm bài viết…',
                    'aria-label': 'Từ khóa tìm kiếm',
                    'autocomplete': 'off',
                  },
                  onInput: _onKeywordChanged,
                ),
                button(
                  type: .button,
                  classes: 'article-search-close',
                  attributes: const {'aria-label': 'Đóng tìm kiếm'},
                  onClick: _close,
                  [.text('×')],
                ),
              ],
            ),
            _buildContent(),
          ],
        ),
      ],
    );
  }

  Component _buildContent() {
    if (_isLoading) {
      return p(
        classes: 'article-search-status',
        attributes: const {'role': 'status'},
        [
          .text('Đang tìm kiếm…'),
        ],
      );
    }
    if (_error case final error?) {
      return p(
        classes: 'article-search-status article-search-status--error',
        attributes: const {'role': 'alert'},
        [
          .text(error),
        ],
      );
    }
    if (_keyword.trim().isEmpty) {
      return p(classes: 'article-search-status', [.text('Nhập từ khóa để tìm bài viết.')]);
    }
    if (_results.isEmpty) {
      return p(
        classes: 'article-search-status',
        attributes: const {'role': 'status'},
        [
          .text('Không tìm thấy bài viết phù hợp.'),
        ],
      );
    }

    return div(
      classes: 'article-search-results',
      attributes: const {'role': 'list'},
      [
        for (final post in _results)
          a(
            href: post.id == null ? '/posts' : '/posts/${post.id}',
            classes: 'article-search-result',
            attributes: const {'role': 'listitem'},
            [
              h2(classes: 'article-search-result-title', [.text(searchResultTitle(post))]),
              if (searchResultExcerpt(post).isNotEmpty)
                p(classes: 'article-search-result-excerpt', [.text(searchResultExcerpt(post))]),
            ],
          ),
      ],
    );
  }
}

String _normalizedBaseUrl(String value) {
  final trimmed = value.trim();
  final absolute = Uri.tryParse(trimmed)?.hasScheme == true
      ? trimmed
      : kIsWeb
      ? web.window.location.origin + (trimmed.startsWith('/') ? trimmed : '/$trimmed')
      : apiBaseUrl;
  return absolute.endsWith('/') ? absolute : '$absolute/';
}

String searchResultTitle(Post post) {
  final title = _plainSearchText(post.title ?? '');
  return title.isEmpty ? 'Bài viết chưa có tiêu đề' : title;
}

String searchResultExcerpt(Post post) {
  final value = _plainSearchText(post.desc?.trim().isNotEmpty == true ? post.desc! : post.content ?? '');
  if (value.length <= 180) return value;
  return '${value.substring(0, 177)}...';
}

String _plainSearchText(String value) {
  return value
      .replaceAll(RegExp(r'<[^>]*>'), '')
      .replaceAll(RegExp(r'[`#>*_\[\]()]'), '')
      .replaceAll(RegExp(r'\s+'), ' ')
      .trim();
}
