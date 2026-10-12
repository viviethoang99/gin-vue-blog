import '../../models/home_page_data.dart';
import '../../models/post_detail.dart' as view;
import '../../models/post_tag.dart';
import '../model/post.dart';

extension PostPresentationMapper on Post {
  PostPreview toPreview({String? href}) {
    return PostPreview(
      title: title?.trim().isNotEmpty == true ? title!.trim() : 'Bài viết chưa có tiêu đề',
      excerpt: _excerpt,
      publishedAt: _formatDate(createdAt),
      category: category?.name.trim().isNotEmpty == true ? category!.name.trim() : 'Chưa phân loại',
      tags: [
        for (final tag in tags ?? const <Category>[])
          if (tag.name.trim().isNotEmpty) PostTag(id: tag.id, name: tag.name.trim()),
      ],
      href: href,
    );
  }

  view.PostDetail toPostDetail({
    required String slug,
    Map<String, int> tagIdsByName = const {},
  }) {
    return view.PostDetail(
      slug: slug,
      title: title?.trim().isNotEmpty == true ? title!.trim() : 'Bài viết chưa có tiêu đề',
      publishedAt: _formatDate(createdAt),
      category: category?.name.trim().isNotEmpty == true ? category!.name.trim() : 'Chưa phân loại',
      tags: [
        for (final tag in tags ?? const <Category>[])
          if (tag.name.trim().isNotEmpty)
            PostTag(
              id: tagIdsByName[tag.name.trim().toLowerCase()] ?? tag.id,
              name: tag.name.trim(),
            ),
      ],
      contentMarkdown: content ?? '',
    );
  }

  String get _excerpt {
    final description = desc?.trim();
    if (description != null && description.isNotEmpty) return description;

    final plainContent = (content ?? '')
        .replaceAll(RegExp(r'[`#>*_\[\]()]'), '')
        .replaceAll(RegExp(r'\s+'), ' ')
        .trim();
    if (plainContent.length <= 180) return plainContent;
    return '${plainContent.substring(0, 177)}...';
  }
}

String _formatDate(DateTime? value) {
  if (value == null) return '';
  final local = value.toLocal();
  String twoDigits(int number) => number.toString().padLeft(2, '0');
  return '${twoDigits(local.day)}/${twoDigits(local.month)}/${local.year}';
}
