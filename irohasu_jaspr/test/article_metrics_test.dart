import 'package:flutter_test/flutter_test.dart';
import 'package:irohasu_jaspr/services/article_metrics.dart';

void main() {
  test('countArticleWords ignores common Markdown and HTML markup', () {
    expect(
      countArticleWords('## Xin chào\nĐây là **một bài** <span>Jaspr</span>.'),
      7,
    );
  });

  test('formatVietnameseInteger adds thousand separators', () {
    expect(formatVietnameseInteger(1800), '1.800');
    expect(formatVietnameseInteger(42), '42');
  });
}
