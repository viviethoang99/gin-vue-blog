import 'post_tag.dart';

class PostDetail {
  const PostDetail({
    required this.slug,
    required this.title,
    required this.publishedAt,
    required this.category,
    required this.tags,
    required this.contentMarkdown,
  });

  final String slug;
  final String title;
  final String publishedAt;
  final String category;
  final List<PostTag> tags;
  final String contentMarkdown;
}
