import 'package:json_annotation/json_annotation.dart';

part 'gist.g.dart';

@JsonSerializable(fieldRename: FieldRename.snake, checked: true)
class Gist {
  final String? url;
  final String? forksUrl;
  final String? commitsUrl;
  final String? id;
  final String? nodeId;
  final String? gitPullUrl;
  final String? gitPushUrl;
  final String? htmlUrl;
  final Map<String, FileValue>? files;
  final bool? public;
  final DateTime? createdAt;
  final DateTime? updatedAt;
  final String? description;
  final int? comments;
  final String? commentsUrl;
  final bool? truncated;

  const Gist({
    this.url,
    this.forksUrl,
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

  factory Gist.fromJson(Map<String, dynamic> json) => _$GistFromJson(json);

  Map<String, dynamic> toJson() => _$GistToJson(this);
}

@JsonSerializable(fieldRename: FieldRename.snake, checked: true)
class FileValue {
  final String? filename;
  final String? type;
  final String? language;
  final String? rawUrl;
  final int? size;

  const FileValue({
    required this.filename,
    required this.type,
    required this.language,
    required this.rawUrl,
    required this.size,
  });

  factory FileValue.fromJson(Map<String, dynamic> json) => _$FileValueFromJson(json);

  Map<String, dynamic> toJson() => _$FileValueToJson(this);
}
