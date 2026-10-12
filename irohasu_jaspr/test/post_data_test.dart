import 'package:flutter_test/flutter_test.dart';
import 'package:irohasu_jaspr/data/post_api/api_result.dart';
import 'package:irohasu_jaspr/data/post_api/list_posts.dart';
import 'package:irohasu_jaspr/data/post_api/post_mapper.dart';
import 'package:irohasu_jaspr/data/model/post.dart';
import 'package:irohasu_jaspr/models/home_page_data.dart';

void main() {
  const categoryJson = {
    'id': 2,
    'created_at': '2026-09-08T00:00:00.000Z',
    'updated_at': '2026-09-08T00:00:00.000Z',
    'name': 'Jaspr',
  };

  test('deserializes the article-list API envelope and snake_case fields', () {
    final response = ListPosts.fromJson({
      'code': 0,
      'message': 'OK',
      'data': [
        {
          'id': 1,
          'created_at': '2026-09-08T00:00:00.000Z',
          'title': 'Irohasu',
          'desc': 'A Dart blog',
          'is_top': true,
          'category': categoryJson,
          'tags': [categoryJson],
        },
      ],
    });

    expect(response.code, 0);
    expect(response.message, 'OK');
    expect(response.data, hasLength(1));
    expect(response.data!.single.isTop, isTrue);
    expect(response.data!.single.category!.name, 'Jaspr');
    final preview = response.data!.single.toPreview();
    expect(preview.tags.single.id, 2);
    expect(preview.tags.single.name, 'Jaspr');
    expect(preview.tags.single.href, '/tags/2');
    expect(PostPreview.fromMap(preview.toMap()).tags.single.id, 2);
    expect(
      response.data!.single
          .toPostDetail(
            slug: 'irohasu',
            tagIdsByName: const {'jaspr': 99},
          )
          .tags
          .single
          .href,
      '/tags/99',
    );
  });

  test('deserializes the paginated article-list response used by the home page', () {
    final response = ListPosts.fromJson({
      'code': 0,
      'message': 'OK',
      'data': {
        'page_num': 1,
        'page_size': 10,
        'total': 1,
        'page_data': [
          {
            'id': 1,
            'created_at': '2026-09-08T00:00:00.000Z',
            'title': 'Bài viết gần đây',
            'category': categoryJson,
          },
        ],
      },
    });

    expect(response.data, hasLength(1));
    expect(response.data!.single.title, 'Bài viết gần đây');
    expect(response.paging?.page, 1);
    expect(response.paging?.limit, 10);
    expect(response.paging?.total, 1);
  });

  test('deserializes a generic article-detail API envelope', () {
    final response = ApiResult<Post>.fromJson(
      {
        'code': 0,
        'message': 'OK',
        'data': {
          'id': 1,
          'created_at': '2026-09-08T00:00:00.000Z',
          'title': 'Irohasu',
          'content': '# Markdown',
          'category': categoryJson,
        },
      },
      (json) => Post.fromJson(json! as Map<String, dynamic>),
    );

    expect(response.data.content, '# Markdown');
    expect(response.data.toPostDetail(slug: 'irohasu').contentMarkdown, '# Markdown');
  });
}
