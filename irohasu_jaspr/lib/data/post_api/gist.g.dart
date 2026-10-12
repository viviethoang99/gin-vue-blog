// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'gist.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Gist _$GistFromJson(Map<String, dynamic> json) => $checkedCreate(
  'Gist',
  json,
  ($checkedConvert) {
    final val = Gist(
      url: $checkedConvert('url', (v) => v as String?),
      forksUrl: $checkedConvert('forks_url', (v) => v as String?),
      commitsUrl: $checkedConvert('commits_url', (v) => v as String?),
      id: $checkedConvert('id', (v) => v as String?),
      nodeId: $checkedConvert('node_id', (v) => v as String?),
      gitPullUrl: $checkedConvert('git_pull_url', (v) => v as String?),
      gitPushUrl: $checkedConvert('git_push_url', (v) => v as String?),
      htmlUrl: $checkedConvert('html_url', (v) => v as String?),
      files: $checkedConvert(
        'files',
        (v) => (v as Map<String, dynamic>?)?.map(
          (k, e) => MapEntry(k, FileValue.fromJson(e as Map<String, dynamic>)),
        ),
      ),
      public: $checkedConvert('public', (v) => v as bool?),
      createdAt: $checkedConvert(
        'created_at',
        (v) => v == null ? null : DateTime.parse(v as String),
      ),
      updatedAt: $checkedConvert(
        'updated_at',
        (v) => v == null ? null : DateTime.parse(v as String),
      ),
      description: $checkedConvert('description', (v) => v as String?),
      comments: $checkedConvert('comments', (v) => (v as num?)?.toInt()),
      commentsUrl: $checkedConvert('comments_url', (v) => v as String?),
      truncated: $checkedConvert('truncated', (v) => v as bool?),
    );
    return val;
  },
  fieldKeyMap: const {
    'forksUrl': 'forks_url',
    'commitsUrl': 'commits_url',
    'nodeId': 'node_id',
    'gitPullUrl': 'git_pull_url',
    'gitPushUrl': 'git_push_url',
    'htmlUrl': 'html_url',
    'createdAt': 'created_at',
    'updatedAt': 'updated_at',
    'commentsUrl': 'comments_url',
  },
);

Map<String, dynamic> _$GistToJson(Gist instance) => <String, dynamic>{
  'url': instance.url,
  'forks_url': instance.forksUrl,
  'commits_url': instance.commitsUrl,
  'id': instance.id,
  'node_id': instance.nodeId,
  'git_pull_url': instance.gitPullUrl,
  'git_push_url': instance.gitPushUrl,
  'html_url': instance.htmlUrl,
  'files': instance.files,
  'public': instance.public,
  'created_at': instance.createdAt?.toIso8601String(),
  'updated_at': instance.updatedAt?.toIso8601String(),
  'description': instance.description,
  'comments': instance.comments,
  'comments_url': instance.commentsUrl,
  'truncated': instance.truncated,
};

FileValue _$FileValueFromJson(Map<String, dynamic> json) =>
    $checkedCreate('FileValue', json, ($checkedConvert) {
      final val = FileValue(
        filename: $checkedConvert('filename', (v) => v as String?),
        type: $checkedConvert('type', (v) => v as String?),
        language: $checkedConvert('language', (v) => v as String?),
        rawUrl: $checkedConvert('raw_url', (v) => v as String?),
        size: $checkedConvert('size', (v) => (v as num?)?.toInt()),
      );
      return val;
    }, fieldKeyMap: const {'rawUrl': 'raw_url'});

Map<String, dynamic> _$FileValueToJson(FileValue instance) => <String, dynamic>{
  'filename': instance.filename,
  'type': instance.type,
  'language': instance.language,
  'raw_url': instance.rawUrl,
  'size': instance.size,
};
