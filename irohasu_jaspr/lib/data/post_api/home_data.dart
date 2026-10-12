import 'package:json_annotation/json_annotation.dart';

part 'home_data.g.dart';

@JsonSerializable(fieldRename: FieldRename.snake, checked: true)
class HomeData {
  const HomeData({
    this.articleCount = 0,
    this.userCount = 0,
    this.messageCount = 0,
    this.categoryCount = 0,
    this.tagCount = 0,
    this.viewCount = 0,
    required this.blogConfig,
  });

  final int articleCount;
  final int userCount;
  final int messageCount;
  final int categoryCount;
  final int tagCount;
  final int viewCount;
  final BlogConfig blogConfig;

  factory HomeData.fromJson(Map<String, dynamic> json) => _$HomeDataFromJson(json);

  Map<String, dynamic> toJson() => _$HomeDataToJson(this);
}

@JsonSerializable(fieldRename: FieldRename.snake, checked: true)
class BlogConfig {
  const BlogConfig({
    this.about = '',
    this.gitee = '',
    this.github = '',
    this.userAvatar = '',
    this.websiteAuthor = '',
    this.websiteAvatar = '',
    this.websiteCreatetime = '',
    this.websiteIntro = '',
    this.websiteName = '',
    this.websiteNotice = '',
  });

  final String about;
  final String gitee;
  final String github;
  final String userAvatar;
  final String websiteAuthor;
  final String websiteAvatar;
  final String websiteCreatetime;
  final String websiteIntro;
  final String websiteName;
  final String websiteNotice;

  factory BlogConfig.fromJson(Map<String, dynamic> json) => _$BlogConfigFromJson(json);

  Map<String, dynamic> toJson() => _$BlogConfigToJson(this);
}
