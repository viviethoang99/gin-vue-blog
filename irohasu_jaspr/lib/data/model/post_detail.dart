import 'package:json_annotation/json_annotation.dart';

part 'post_detail.g.dart';

@JsonSerializable(fieldRename: FieldRename.snake, checked: true)
class PostDetail {
  final String? url;
  final String? commitsUrl;
  final String? id;
  final String? nodeId;
  final String? gitPullUrl;
  final String? gitPushUrl;
  final String? htmlUrl;
  @JsonKey(name: 'files', fromJson: _fromListFiles)
  final List<GistFile>? files;
  final bool? public;
  final DateTime? createdAt;
  final DateTime? updatedAt;
  final String? description;
  final int? comments;
  final String? commentsUrl;
  final bool? truncated;

  const PostDetail({
    this.url,
    this.commitsUrl,
    this.id,
    this.nodeId,
    this.gitPullUrl,
    this.gitPushUrl,
    this.htmlUrl,
    this.files,
    this.public,
    this.createdAt,
    this.updatedAt,
    this.description,
    this.comments,
    this.commentsUrl,
    this.truncated,
  });

  static List<GistFile> _fromListFiles(Map<String, dynamic> value) {
    return value.entries.map((e) => GistFile.fromJson(e.value as Map<String, dynamic>)).toList();
  }

  factory PostDetail.fromJson(Map<String, dynamic> json) => _$PostDetailFromJson(json);

  Map<String, dynamic> toJson() => _$PostDetailToJson(this);
}

@JsonSerializable(fieldRename: FieldRename.snake, checked: true)
class GistFile {
  final String? filename;
  final String? type;
  final String? language;
  final String? rawUrl;
  final int? size;
  final bool? truncated;
  final String? content;

  const GistFile({
    this.filename,
    this.type,
    this.language,
    this.rawUrl,
    this.size,
    this.truncated,
    this.content,
  });

  factory GistFile.fromJson(Map<String, dynamic> json) => _$GistFileFromJson(json);

  Map<String, dynamic> toJson() => _$GistFileToJson(this);
}
