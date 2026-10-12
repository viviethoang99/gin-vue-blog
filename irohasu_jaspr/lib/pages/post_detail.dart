import 'dart:async';

import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';
import 'package:universal_web/js_interop.dart';
import 'package:universal_web/web.dart' as web;

import '../components/markdown_content.dart';
import '../config/environment.dart';
import '../constants/theme.dart';
import '../models/post_detail.dart';

class PostDetailPage extends StatelessComponent {
  const PostDetailPage({required this.post, super.key});

  final PostDetail post;

  @override
  Component build(BuildContext context) {
    final headings = extractMarkdownHeadings(post.contentMarkdown);

    return section(classes: 'post-detail-page', [
      main_(classes: 'post-detail-shell', [
        article(classes: 'post-detail-content', [
          header(classes: 'post-detail-header', [
            h1(classes: 'post-detail-title', [.text(post.title)]),
            _PostDetailMeta(post: post),
          ]),
          MarkdownContent(
            source: post.contentMarkdown,
            imageBaseUrl: publicApiBaseUrl,
          ),
        ]),
        if (headings.isNotEmpty)
          TableOfContents(
            headings: [
              for (final heading in headings) {'level': heading.level, 'text': heading.text, 'id': heading.id},
            ],
          ),
      ]),
    ]);
  }

  @css
  static List<StyleRule> get styles => [
    css('.post-detail-page', [
      css('&').styles(
        width: 100.percent,
        minHeight: 100.vh,
        padding: .only(top: 118.px, right: 80.px, bottom: 90.px, left: 80.px),
        margin: .zero,
        boxSizing: .borderBox,
        color: textSecondaryColor,
        backgroundColor: Colors.transparent,
      ),
      css('.post-detail-shell').styles(
        display: .flex,
        width: 100.percent,
        maxWidth: 1080.px,
        margin: .only(left: Unit.expression('calc((100vw - 720px) / 2 - 80px)'), right: .zero),
        alignItems: .start,
        gap: .all(38.px),
      ),
      css('.post-detail-content').styles(
        width: 100.percent,
        minWidth: 0.px,
        maxWidth: 720.px,
        flex: .none,
      ),
      css('.post-detail-header').styles(margin: .only(bottom: 36.px)),
      css('.post-detail-title').styles(
        color: textPrimaryColor,
        textAlign: .center,
        fontSize: 1.65.rem,
        fontWeight: .w600,
        lineHeight: 1.35.em,
      ),
      css('.post-detail-meta').styles(
        display: .flex,
        margin: .only(top: 6.px),
        flexWrap: .wrap,
        justifyContent: .center,
        gap: Gap(row: 3.px, column: 7.px),
        color: primaryColor,
        fontSize: 0.76.rem,
        lineHeight: 1.4.em,
      ),
      css('.post-detail-meta-item:not(:last-child)::after').styles(
        content: '·',
        margin: .only(left: 7.px),
        color: textMutedColor,
      ),
      css('.post-detail-meta a').styles(color: .inherit, textDecoration: .none),
      css('.post-detail-meta a:hover, .post-detail-meta a:focus-visible').styles(
        textDecoration: TextDecoration(line: .underline),
      ),
      css('.markdown-body').styles(
        color: textSecondaryColor,
        fontSize: 0.98.rem,
        lineHeight: 1.72.em,
        raw: const {'word-break': 'break-word'},
      ),
      css('.markdown-body a', [
        css('&').styles(
          color: primaryColor,
          textDecoration: TextDecoration(line: .none),
        ),
        css('&:hover').styles(textDecoration: TextDecoration(line: .underline)),
      ]),
      css('.markdown-body p').styles(margin: .only(bottom: 16.px)),
      css('.markdown-body h2').styles(
        position: .relative(),
        padding: .only(bottom: 8.px),
        margin: .only(top: 30.px, bottom: 16.px),
        border: .only(
          bottom: BorderSide.dashed(color: borderColor, width: 1.px),
        ),
        color: textPrimaryColor,
        fontSize: 1.35.rem,
        fontWeight: .w600,
        lineHeight: 1.4.em,
        raw: const {'scroll-margin-top': '88px'},
      ),
      css('.markdown-body h2::after').styles(
        content: '',
        display: .block,
        position: .absolute(bottom: (-1).px, left: 0.px),
        width: 34.px,
        height: 2.px,
        backgroundColor: primaryColor,
      ),
      css('.markdown-body h3').styles(
        margin: .only(top: 24.px, bottom: 10.px),
        color: textPrimaryColor,
        fontSize: 1.08.rem,
        fontWeight: .w600,
        lineHeight: 1.45.em,
        raw: const {'scroll-margin-top': '88px'},
      ),
      css('.markdown-body blockquote').styles(
        padding: .symmetric(vertical: 16.px, horizontal: 22.px),
        margin: .symmetric(vertical: 20.px),
        border: .only(
          left: BorderSide.solid(color: noteTextColor, width: 3.px),
        ),
        color: noteTextColor,
        backgroundColor: noteBackgroundColor,
      ),
      css('.markdown-body blockquote p:last-child').styles(margin: .zero),
      css('.markdown-body ul, .markdown-body ol').styles(
        padding: .only(left: 26.px),
        margin: .only(top: 8.px, bottom: 20.px),
      ),
      css('.markdown-body li').styles(margin: .only(bottom: 6.px)),
      css('.markdown-body code').styles(
        padding: .symmetric(vertical: 2.px, horizontal: 5.px),
        radius: .all(.circular(4.px)),
        color: textPrimaryColor,
        fontFamily: const .list([FontFamily('SFMono-Regular'), FontFamily('Consolas'), FontFamilies.monospace]),
        fontSize: 0.86.em,
        backgroundColor: surfaceColor,
      ),
      css('.markdown-body pre').styles(
        padding: .all(18.px),
        margin: .only(top: 16.px, bottom: 20.px),
        border: .all(style: .solid, color: borderColor, width: 1.px),
        radius: .all(.circular(8.px)),
        overflow: .auto,
        backgroundColor: surfaceColor,
      ),
      css('.markdown-body pre code').styles(padding: .zero, backgroundColor: Colors.transparent),
      css('.markdown-body hr').styles(
        margin: .symmetric(vertical: 32.px),
        border: .only(
          top: BorderSide.solid(color: borderColor, width: 1.px),
        ),
      ),
      css('.markdown-body p:has(> img)').styles(
        display: .flex,
        flexWrap: .wrap,
        gap: .all(8.px),
      ),
      css('.markdown-body img').styles(
        display: .block,
        width: Unit.expression('calc(50% - 4px)'),
        height: .auto,
        maxWidth: 100.percent,
        radius: .all(.circular(5.px)),
        raw: const {'object-fit': 'cover'},
      ),
      css('.markdown-body p > img:only-child').styles(width: 100.percent),
      css('.markdown-body table').styles(
        width: 100.percent,
        margin: .symmetric(vertical: 20.px),
        border: .all(style: .solid, color: borderColor, width: 1.px),
        raw: const {'border-collapse': 'collapse'},
      ),
      css('.markdown-body th, .markdown-body td').styles(
        padding: .symmetric(vertical: 8.px, horizontal: 10.px),
        border: .all(style: .solid, color: borderColor, width: 1.px),
        textAlign: .left,
      ),
      css('.article-toc').styles(
        position: .sticky(top: 190.px),
        width: 220.px,
        opacity: 0.48,
        transition: Transition('opacity', duration: 180.ms, curve: .easeOut),
        flex: .none,
      ),
      css('.article-toc:hover').styles(opacity: 0.9),
      css('.article-toc ul').styles(
        padding: .zero,
        margin: .zero,
        listStyle: .none,
      ),
      css('.article-toc li').styles(margin: .only(bottom: 8.px)),
      css('.article-toc li.toc-level-3').styles(padding: .only(left: 14.px)),
      css('.article-toc a').styles(
        display: .block,
        overflow: .hidden,
        color: textMutedColor,
        fontSize: 0.78.rem,
        textDecoration: TextDecoration(line: .none),
        textOverflow: .ellipsis,
        whiteSpace: .noWrap,
      ),
      css('.article-toc a::before').styles(
        content: '–',
        margin: .only(right: 7.px),
      ),
      css('.article-toc a:hover').styles(color: primaryColor),
      css('.article-toc a.active').styles(
        color: primaryColor,
        fontWeight: .w600,
        textDecoration: TextDecoration(line: .underline),
        raw: const {'text-underline-offset': '4px'},
      ),
    ]),
    css.media(MediaQuery.screen(maxWidth: 1280.px), [
      css('.post-detail-page .post-detail-shell').styles(
        maxWidth: 720.px,
        margin: .only(left: .auto, right: .auto),
      ),
      css('.post-detail-page .article-toc').styles(display: .none),
    ]),
    css.media(MediaQuery.screen(maxWidth: 700.px), [
      css('.post-detail-page').styles(
        padding: .only(top: 104.px, right: 28.px, bottom: 60.px, left: 28.px),
      ),
      css('.post-detail-page .post-detail-title').styles(fontSize: 1.42.rem),
      css('.post-detail-page .markdown-body').styles(fontSize: 0.94.rem, lineHeight: 1.68.em),
      css('.post-detail-page .markdown-body h2').styles(fontSize: 1.22.rem),
      css('.post-detail-page .markdown-body img').styles(width: 100.percent),
    ]),
    css.media(const MediaQuery.raw('(prefers-reduced-motion: reduce)'), [
      css('.post-detail-page .article-toc').styles(raw: const {'transition': 'none'}),
    ]),
  ];
}

class _PostDetailMeta extends StatelessComponent {
  const _PostDetailMeta({required this.post});

  final PostDetail post;

  @override
  Component build(BuildContext context) {
    return div(classes: 'post-detail-meta', [
      span(classes: 'post-detail-meta-item', [.text(post.publishedAt)]),
      span(classes: 'post-detail-meta-item', [.text(post.category)]),
      for (final tag in post.tags)
        span(classes: 'post-detail-meta-item', [
          a(href: tag.href, [.text('#${tag.name}')]),
        ]),
    ]);
  }
}

@client
class TableOfContents extends StatefulComponent {
  const TableOfContents({required this.headings, super.key});

  final List<Map<String, dynamic>> headings;

  @override
  State<TableOfContents> createState() => _TableOfContentsState();
}

class _TableOfContentsState extends State<TableOfContents> {
  StreamSubscription<web.Event>? _scrollSubscription;
  Timer? _setupTimer;
  int? _animationFrameId;
  String? _activeId;
  bool _disposed = false;

  @override
  void initState() {
    super.initState();
    _activeId = _headingId(component.headings.isEmpty ? null : component.headings.first);
    if (!kIsWeb) return;

    _setupTimer = Timer(Duration.zero, () {
      if (_disposed) return;
      _updateActiveHeading();
      _scrollSubscription = web.EventStreamProviders.scrollEvent.forTarget(web.window).listen((_) {
        _scheduleActiveHeadingUpdate();
      });
    });
  }

  void _scheduleActiveHeadingUpdate() {
    if (_animationFrameId != null || _disposed) return;
    _animationFrameId = web.window.requestAnimationFrame(_onAnimationFrame.toJS);
  }

  void _onAnimationFrame(double _) {
    _animationFrameId = null;
    _updateActiveHeading();
  }

  void _updateActiveHeading() {
    if (!kIsWeb || _disposed) return;

    String? nextId = _headingId(component.headings.isEmpty ? null : component.headings.first);
    const activationOffset = 120.0;
    for (final heading in component.headings) {
      final id = _headingId(heading);
      if (id == null) continue;
      final element = web.document.getElementById(id);
      if (element == null) continue;
      if (element.getBoundingClientRect().top <= activationOffset) {
        nextId = id;
      } else {
        break;
      }
    }

    if (nextId == null || nextId == _activeId) return;
    setState(() => _activeId = nextId);
  }

  void _activateHeading(String id) {
    if (id != _activeId) setState(() => _activeId = id);
    _scrollToHeading(id);
  }

  @override
  Component build(BuildContext context) {
    return aside(classes: 'article-toc', [
      nav(
        attributes: const {'aria-label': 'Mục lục bài viết'},
        [
          ul([
            for (final heading in component.headings)
              li(classes: 'toc-level-${heading['level']}', [
                a(
                  href: '#${Uri.encodeComponent(_headingId(heading)!)}',
                  classes: _headingId(heading) == _activeId ? 'active' : null,
                  attributes: {
                    if (_headingId(heading) == _activeId) 'aria-current': 'location',
                  },
                  events: {
                    'click': (event) {
                      event.preventDefault();
                      _activateHeading(_headingId(heading)!);
                    },
                  },
                  [.text(heading['text']! as String)],
                ),
              ]),
          ]),
        ],
      ),
    ]);
  }

  @override
  void dispose() {
    _disposed = true;
    _setupTimer?.cancel();
    _scrollSubscription?.cancel();
    final animationFrameId = _animationFrameId;
    if (animationFrameId != null && kIsWeb) web.window.cancelAnimationFrame(animationFrameId);
    super.dispose();
  }
}

String? _headingId(Map<String, dynamic>? heading) => heading?['id'] as String?;

void _scrollToHeading(String id) {
  if (!kIsWeb) return;
  final target = web.document.getElementById(id);
  if (target == null) return;

  final reduceMotion = web.window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  target.scrollIntoView(
    web.ScrollIntoViewOptions(
      behavior: reduceMotion ? 'auto' : 'smooth',
      block: 'start',
    ),
  );
  web.window.history.replaceState(null, '', '#${Uri.encodeComponent(id)}');
}
