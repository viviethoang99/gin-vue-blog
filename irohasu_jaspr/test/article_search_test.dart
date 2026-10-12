import 'package:flutter_test/flutter_test.dart';
import 'package:irohasu_jaspr/components/article_search.dart';
import 'package:irohasu_jaspr/data/model/post.dart';

void main() {
  group('article search presentation', () {
    test('strips backend highlight markup from result excerpts', () {
      const post = Post(
        content: "Xây dựng với <span style='color:red'>Vue</span> 3\n\n- Pinia",
      );

      expect(searchResultExcerpt(post), 'Xây dựng với Vue 3 - Pinia');
    });

    test('uses a safe fallback for an empty title', () {
      const post = Post(title: '   ');

      expect(searchResultTitle(post), 'Bài viết chưa có tiêu đề');
    });

    test('limits long excerpts to 180 characters', () {
      final post = Post(content: List.filled(200, 'a').join());

      expect(searchResultExcerpt(post), hasLength(180));
      expect(searchResultExcerpt(post), endsWith('...'));
    });
  });
}
