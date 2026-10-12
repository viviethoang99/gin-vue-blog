import '../../models/home_page_data.dart';
import 'home_data.dart';

extension HomePresentationMapper on HomeData {
  HomePageData toPageData({
    required List<PostPreview> recentPosts,
  }) {
    final config = blogConfig;
    final biography = [
      config.websiteNotice.trim(),
      _plainText(config.about),
    ].where((value) => value.isNotEmpty).map((value) => RichParagraph([TextSegment(value)])).toList(growable: false);

    final links = <FooterLink>[
      const FooterLink(label: 'Home', href: '/'),
      if (_isHttpUrl(config.github)) FooterLink(label: 'GitHub', href: config.github.trim()),
      if (_isHttpUrl(config.gitee)) FooterLink(label: 'Gitee', href: config.gitee.trim()),
    ];

    return HomePageData(
      profile: AuthorProfile(
        name: config.websiteAuthor.trim(),
        tagline: config.websiteIntro.trim(),
        avatarUrl: config.websiteAvatar.trim(),
        biography: biography,
      ),
      recentPosts: recentPosts,
      footer: SiteFooterData(
        owner: config.websiteAuthor.trim(),
        startYear: _yearFrom(config.websiteCreatetime),
        links: links,
      ),
    );
  }
}

int? _yearFrom(String value) {
  final match = RegExp(r'^\s*(\d{4})').firstMatch(value);
  return int.tryParse(match?.group(1) ?? '');
}

bool _isHttpUrl(String value) {
  final uri = Uri.tryParse(value.trim());
  return uri != null && (uri.scheme == 'http' || uri.scheme == 'https') && uri.host.isNotEmpty;
}

String _plainText(String markdown) => markdown
    .replaceAll(RegExp(r'```[^\n]*\n?'), '')
    .replaceAll('```', '')
    .replaceAll(RegExp(r'[#>*_`]'), '')
    .replaceAll(RegExp(r'\s+'), ' ')
    .trim();
