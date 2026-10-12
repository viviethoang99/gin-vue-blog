import 'package:flutter_test/flutter_test.dart';
import 'package:irohasu_jaspr/data/post_api/api_result.dart';
import 'package:irohasu_jaspr/data/post_api/home_data.dart';
import 'package:irohasu_jaspr/data/post_api/home_mapper.dart';
import 'package:irohasu_jaspr/models/home_page_data.dart';

void main() {
  const responseJson = <String, dynamic>{
    'code': 0,
    'message': 'OK',
    'data': <String, dynamic>{
      'article_count': 4,
      'user_count': 5,
      'message_count': 0,
      'category_count': 3,
      'tag_count': 2,
      'view_count': 10,
      'blog_config': <String, dynamic>{
        'about': '```dart\nprint("hello");\n```',
        'github': 'https://github.com/example',
        'gitee': '',
        'website_author': 'Remote author',
        'website_avatar': 'https://example.com/avatar.png',
        'website_createtime': '2024-02-20 22:40:22',
        'website_intro': 'Remote intro',
        'website_name': 'Remote blog',
        'website_notice': 'Remote notice',
      },
    },
  };

  test('decodes and maps the home endpoint response for presentation', () {
    final response = ApiResult<HomeData>.fromJson(
      responseJson,
      (value) => HomeData.fromJson(value! as Map<String, dynamic>),
    );
    final page = response.data.toPageData(
      recentPosts: const <PostPreview>[],
    );

    expect(response.code, 0);
    expect(page.profile.name, 'Remote author');
    expect(page.profile.tagline, 'Remote intro');
    expect(page.profile.avatarUrl, 'https://example.com/avatar.png');
    expect(page.profile.biography.first.segments.single.text, 'Remote notice');
    expect(page.profile.biography.last.segments.single.text, 'print("hello");');
    expect(page.footer.startYear, 2024);
    expect(page.footer.links.map((link) => link.label), contains('GitHub'));
  });

  test('keeps missing home API fields empty so the UI can hide them', () {
    const data = HomeData(blogConfig: BlogConfig());

    final page = data.toPageData(recentPosts: const <PostPreview>[]);

    expect(page.profile.name, isEmpty);
    expect(page.profile.tagline, isEmpty);
    expect(page.profile.avatarUrl, isEmpty);
    expect(page.profile.biography, isEmpty);
    expect(page.footer.owner, isEmpty);
    expect(page.footer.startYear, isNull);
    expect(page.footer.links.map((link) => link.label), ['Home']);
  });
}
