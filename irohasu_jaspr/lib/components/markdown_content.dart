import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';
import 'package:markdown/markdown.dart' as markdown;

class MarkdownHeading {
  const MarkdownHeading({required this.level, required this.text, required this.id});

  final int level;
  final String text;
  final String id;
}

List<MarkdownHeading> extractMarkdownHeadings(String source) {
  final nodes = _parse(source);
  final headings = <MarkdownHeading>[];
  final slugCounts = <String, int>{};

  void visit(markdown.Node node) {
    if (node case final markdown.Element element) {
      final level = switch (element.tag) {
        'h2' => 2,
        'h3' => 3,
        _ => null,
      };
      if (level != null) {
        final text = element.textContent;
        headings.add(MarkdownHeading(level: level, text: text, id: _uniqueSlug(text, slugCounts)));
      }
      for (final child in element.children ?? const <markdown.Node>[]) {
        visit(child);
      }
    }
  }

  for (final node in nodes) {
    visit(node);
  }
  return headings;
}

class MarkdownContent extends StatelessComponent {
  const MarkdownContent({required this.source, this.imageBaseUrl, super.key});

  final String source;
  final String? imageBaseUrl;

  @override
  Component build(BuildContext context) {
    final renderer = _MarkdownRenderer(imageBaseUrl: imageBaseUrl);
    return div(classes: 'markdown-body', renderer.render(_parse(source)));
  }
}

List<markdown.Node> _parse(String source) => markdown.Document(
  extensionSet: markdown.ExtensionSet.gitHubWeb,
  encodeHtml: true,
).parse(normalizeMarkdownImages(source));

class _MarkdownRenderer {
  _MarkdownRenderer({this.imageBaseUrl});

  final String? imageBaseUrl;
  final Map<String, int> _slugCounts = {};

  List<Component> render(List<markdown.Node> nodes) => [for (final node in nodes) _renderNode(node)];

  Component _renderNode(markdown.Node node) {
    if (node case final markdown.Text text) return .text(text.text);
    if (node is! markdown.Element) return .empty();

    final children = render(node.children ?? const <markdown.Node>[]);
    return switch (node.tag) {
      'h1' => h1(id: _uniqueSlug(node.textContent, _slugCounts), children),
      'h2' => h2(id: _uniqueSlug(node.textContent, _slugCounts), children),
      'h3' => h3(id: _uniqueSlug(node.textContent, _slugCounts), children),
      'h4' => h4(children),
      'h5' => h5(children),
      'h6' => h6(children),
      'p' => p(children),
      'strong' => strong(children),
      'em' => em(children),
      'code' => code(classes: _safeCodeClass(node.attributes['class']), children),
      'pre' => pre(children),
      'blockquote' => blockquote(children),
      'ul' => ul(children),
      'ol' => ol(children),
      'li' => li(children),
      'hr' => hr(),
      'br' => br(),
      'a' => _renderLink(node, children),
      'img' => _renderImage(node),
      'table' => table(children),
      'thead' => thead(children),
      'tbody' => tbody(children),
      'tr' => tr(children),
      'th' => th(children),
      'td' => td(children),
      _ => .fragment(children),
    };
  }

  Component _renderLink(markdown.Element node, List<Component> children) {
    final href = _safeUrl(node.attributes['href'], allowMail: true);
    return href == null ? span(children) : a(href: href, children);
  }

  Component _renderImage(markdown.Element node) {
    final source = resolveMarkdownImageUrl(node.attributes['src'], baseUrl: imageBaseUrl);
    if (source == null) return .empty();
    final title = node.attributes['title']?.trim();
    return img(
      src: source,
      alt: node.attributes['alt'] ?? '',
      classes: 'markdown-image',
      loading: .lazy,
      attributes: {
        'decoding': 'async',
        if (title != null && title.isNotEmpty) 'title': title,
      },
    );
  }
}

String? resolveMarkdownImageUrl(String? value, {String? baseUrl}) {
  final safeValue = _safeUrl(value);
  if (safeValue == null) return null;

  final uri = Uri.parse(safeValue);
  if (uri.hasScheme || baseUrl == null || baseUrl.trim().isEmpty) return safeValue;
  final baseUri = Uri.tryParse(baseUrl.trim());
  if (baseUri == null || !baseUri.hasScheme || (baseUri.scheme != 'http' && baseUri.scheme != 'https')) return null;
  return baseUri.resolveUri(uri).toString();
}

String normalizeMarkdownImages(String source) {
  final imageTag = RegExp(r'(?:^[ \t]*)?\\?<img\b[^>]*>', caseSensitive: false, multiLine: true);
  final wrapperLine = RegExp(
    r'^[ \t]*\\?</?(?:p|a)\b[^>]*>[ \t]*$',
    caseSensitive: false,
    multiLine: true,
  );

  return source
      .replaceAllMapped(imageTag, (match) {
        final attributes = _parseHtmlAttributes(match.group(0)!);
        final rawSource = _unwrapMarkdownUrl(attributes['src']);
        if (_safeUrl(rawSource) == null) return '';

        final alt = _escapeMarkdownLabel(attributes['alt'] ?? '');
        final title = attributes['title']?.trim();
        final titleSuffix = title == null || title.isEmpty ? '' : ' "${title.replaceAll('"', r'\"')}"';
        return '![$alt]($rawSource$titleSuffix)';
      })
      .replaceAll(wrapperLine, '');
}

Map<String, String> _parseHtmlAttributes(String tag) {
  final attributes = <String, String>{};
  final pattern = RegExp(r'''([a-zA-Z_:][\w:.-]*)\s*=\s*(?:"([^"]*)"|'([^']*)')''');
  for (final match in pattern.allMatches(tag)) {
    attributes[match.group(1)!.toLowerCase()] = match.group(2) ?? match.group(3) ?? '';
  }
  return attributes;
}

String? _unwrapMarkdownUrl(String? value) {
  if (value == null) return null;
  final match = RegExp(r'^\[[^\]]*\]\(([^)]+)\)$').firstMatch(value.trim());
  return match?.group(1)?.trim() ?? value.trim();
}

String _escapeMarkdownLabel(String value) =>
    value.replaceAll(r'\', r'\\').replaceAll('[', r'\[').replaceAll(']', r'\]');

String? _safeUrl(String? value, {bool allowMail = false}) {
  if (value == null || value.trim().isEmpty) return null;
  final uri = Uri.tryParse(value.trim());
  if (uri == null) return null;
  if (!uri.hasScheme || uri.scheme == 'http' || uri.scheme == 'https' || (allowMail && uri.scheme == 'mailto')) {
    return value.trim();
  }
  return null;
}

String? _safeCodeClass(String? value) {
  if (value == null || !RegExp(r'^language-[a-zA-Z0-9_+-]+$').hasMatch(value)) return null;
  return value;
}

String _uniqueSlug(String text, Map<String, int> counts) {
  final base = text
      .trim()
      .toLowerCase()
      .replaceAll(RegExp(r'[^\p{L}\p{N}\s-]', unicode: true), '')
      .replaceAll(RegExp(r'[\s-]+'), '-');
  final safeBase = base.isEmpty ? 'section' : base;
  final count = counts.update(safeBase, (value) => value + 1, ifAbsent: () => 1);
  return count == 1 ? safeBase : '$safeBase-$count';
}
