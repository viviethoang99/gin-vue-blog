import 'post_tag.dart';

class HomePageData {
  const HomePageData({
    required this.profile,
    required this.recentPosts,
    required this.footer,
  });

  final AuthorProfile profile;
  final List<PostPreview> recentPosts;
  final SiteFooterData footer;

  factory HomePageData.withRecentPosts(List<PostPreview> recentPosts) {
    return HomePageData(
      profile: const AuthorProfile.empty(),
      recentPosts: recentPosts,
      footer: const SiteFooterData.empty(),
    );
  }
}

class AuthorProfile {
  const AuthorProfile({
    required this.name,
    required this.tagline,
    required this.avatarUrl,
    required this.biography,
  });

  const AuthorProfile.empty() : name = '', tagline = '', avatarUrl = '', biography = const [];

  final String name;
  final String tagline;
  final String avatarUrl;
  final List<RichParagraph> biography;
}

class RichParagraph {
  const RichParagraph(this.segments);

  final List<TextSegment> segments;
}

class TextSegment {
  const TextSegment(this.text, {this.href});

  final String text;
  final String? href;
}

class PostPreview {
  const PostPreview({
    required this.title,
    required this.excerpt,
    required this.publishedAt,
    required this.category,
    required this.tags,
    this.href,
  });

  final String title;
  final String excerpt;
  final String publishedAt;
  final String category;
  final List<PostTag> tags;
  final String? href;

  factory PostPreview.fromMap(Map<String, dynamic> value) {
    return PostPreview(
      title: value['title'] as String,
      excerpt: value['excerpt'] as String,
      publishedAt: value['publishedAt'] as String,
      category: value['category'] as String,
      tags: [
        for (final tag in value['tags'] as List<dynamic>) PostTag.fromMap((tag as Map).cast<String, dynamic>()),
      ],
      href: value['href'] as String?,
    );
  }

  Map<String, dynamic> toMap() => {
    'title': title,
    'excerpt': excerpt,
    'publishedAt': publishedAt,
    'category': category,
    'tags': [for (final tag in tags) tag.toMap()],
    'href': href,
  };
}

class SiteFooterData {
  const SiteFooterData({
    required this.owner,
    required this.startYear,
    required this.links,
  });

  const SiteFooterData.empty() : owner = '', startYear = null, links = const [];

  final String owner;
  final int? startYear;
  final List<FooterLink> links;
}

class FooterLink {
  const FooterLink({required this.label, required this.href});

  final String label;
  final String href;
}
