import 'package:flutter_test/flutter_test.dart';
import 'package:irohasu_jaspr/components/markdown_content.dart';

void main() {
  group('extractMarkdownHeadings', () {
    test('extracts level two and three headings in document order', () {
      final headings = extractMarkdownHeadings('''
# Tiêu đề bài viết
## Phần đầu
### Chi tiết
## Phần cuối
''');

      expect(headings.map((heading) => heading.level), [2, 3, 2]);
      expect(headings.map((heading) => heading.text), ['Phần đầu', 'Chi tiết', 'Phần cuối']);
      expect(headings.map((heading) => heading.id), ['phần-đầu', 'chi-tiết', 'phần-cuối']);
    });

    test('creates unique anchors for duplicate headings', () {
      final headings = extractMarkdownHeadings('## Lặp lại\n\n## Lặp lại');

      expect(headings.map((heading) => heading.id), ['lặp-lại', 'lặp-lại-2']);
    });
  });

  group('resolveMarkdownImageUrl', () {
    test('resolves an API-relative image URL', () {
      expect(
        resolveMarkdownImageUrl(
          '/public/uploaded/photo.jpg',
          baseUrl: 'https://blog.example.com/',
        ),
        'https://blog.example.com/public/uploaded/photo.jpg',
      );
    });

    test('keeps absolute HTTP URLs unchanged', () {
      expect(
        resolveMarkdownImageUrl(
          'https://cdn.example.com/photo.jpg',
          baseUrl: 'https://blog.example.com/',
        ),
        'https://cdn.example.com/photo.jpg',
      );
    });

    test('rejects unsafe image URLs', () {
      expect(resolveMarkdownImageUrl('javascript:alert(1)'), isNull);
      expect(resolveMarkdownImageUrl('data:text/html,unsafe'), isNull);
    });
  });

  group('normalizeMarkdownImages', () {
    test('converts safe inline HTML images and removes their wrappers', () {
      final result = normalizeMarkdownImages('''
<p align="center">
  <a target="_blank" href="#">
    <img style="display: inline-block" src="https://img.shields.io/badge/Go-1.19-blue" alt="Go" />
  </a>
</p>
''');

      expect(result, contains('![Go](https://img.shields.io/badge/Go-1.19-blue)'));
      expect(result, isNot(contains('<p')));
      expect(result, isNot(contains('<a')));
      expect(result, isNot(contains('style=')));
      expect(result, isNot(contains('    ![Go]')));
    });

    test('supports a Markdown-formatted URL inside an HTML src attribute', () {
      final result = normalizeMarkdownImages(
        r'<img src="[badge](https://img.shields.io/badge/Dart-blue)" />',
      );

      expect(result, ' ![](https://img.shields.io/badge/Dart-blue)'.trimLeft());
    });

    test('drops HTML images with unsafe URLs and event attributes', () {
      final result = normalizeMarkdownImages(
        '<img src="javascript:alert(1)" onerror="alert(2)" />',
      );

      expect(result, isEmpty);
    });
  });
}
