// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'home_data.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

HomeData _$HomeDataFromJson(Map<String, dynamic> json) => $checkedCreate(
  'HomeData',
  json,
  ($checkedConvert) {
    final val = HomeData(
      articleCount: $checkedConvert(
        'article_count',
        (v) => (v as num?)?.toInt() ?? 0,
      ),
      userCount: $checkedConvert(
        'user_count',
        (v) => (v as num?)?.toInt() ?? 0,
      ),
      messageCount: $checkedConvert(
        'message_count',
        (v) => (v as num?)?.toInt() ?? 0,
      ),
      categoryCount: $checkedConvert(
        'category_count',
        (v) => (v as num?)?.toInt() ?? 0,
      ),
      tagCount: $checkedConvert('tag_count', (v) => (v as num?)?.toInt() ?? 0),
      viewCount: $checkedConvert(
        'view_count',
        (v) => (v as num?)?.toInt() ?? 0,
      ),
      blogConfig: $checkedConvert(
        'blog_config',
        (v) => BlogConfig.fromJson(v as Map<String, dynamic>),
      ),
    );
    return val;
  },
  fieldKeyMap: const {
    'articleCount': 'article_count',
    'userCount': 'user_count',
    'messageCount': 'message_count',
    'categoryCount': 'category_count',
    'tagCount': 'tag_count',
    'viewCount': 'view_count',
    'blogConfig': 'blog_config',
  },
);

Map<String, dynamic> _$HomeDataToJson(HomeData instance) => <String, dynamic>{
  'article_count': instance.articleCount,
  'user_count': instance.userCount,
  'message_count': instance.messageCount,
  'category_count': instance.categoryCount,
  'tag_count': instance.tagCount,
  'view_count': instance.viewCount,
  'blog_config': instance.blogConfig,
};

BlogConfig _$BlogConfigFromJson(Map<String, dynamic> json) => $checkedCreate(
  'BlogConfig',
  json,
  ($checkedConvert) {
    final val = BlogConfig(
      about: $checkedConvert('about', (v) => v as String? ?? ''),
      gitee: $checkedConvert('gitee', (v) => v as String? ?? ''),
      github: $checkedConvert('github', (v) => v as String? ?? ''),
      userAvatar: $checkedConvert('user_avatar', (v) => v as String? ?? ''),
      websiteAuthor: $checkedConvert(
        'website_author',
        (v) => v as String? ?? '',
      ),
      websiteAvatar: $checkedConvert(
        'website_avatar',
        (v) => v as String? ?? '',
      ),
      websiteCreatetime: $checkedConvert(
        'website_createtime',
        (v) => v as String? ?? '',
      ),
      websiteIntro: $checkedConvert('website_intro', (v) => v as String? ?? ''),
      websiteName: $checkedConvert('website_name', (v) => v as String? ?? ''),
      websiteNotice: $checkedConvert(
        'website_notice',
        (v) => v as String? ?? '',
      ),
    );
    return val;
  },
  fieldKeyMap: const {
    'userAvatar': 'user_avatar',
    'websiteAuthor': 'website_author',
    'websiteAvatar': 'website_avatar',
    'websiteCreatetime': 'website_createtime',
    'websiteIntro': 'website_intro',
    'websiteName': 'website_name',
    'websiteNotice': 'website_notice',
  },
);

Map<String, dynamic> _$BlogConfigToJson(BlogConfig instance) =>
    <String, dynamic>{
      'about': instance.about,
      'gitee': instance.gitee,
      'github': instance.github,
      'user_avatar': instance.userAvatar,
      'website_author': instance.websiteAuthor,
      'website_avatar': instance.websiteAvatar,
      'website_createtime': instance.websiteCreatetime,
      'website_intro': instance.websiteIntro,
      'website_name': instance.websiteName,
      'website_notice': instance.websiteNotice,
    };
