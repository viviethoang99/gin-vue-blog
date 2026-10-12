// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'post.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Post _$PostFromJson(Map<String, dynamic> json) => $checkedCreate(
  'Post',
  json,
  ($checkedConvert) {
    final val = Post(
      id: $checkedConvert('id', (v) => (v as num?)?.toInt()),
      createdAt: $checkedConvert(
        'created_at',
        (v) => v == null ? null : DateTime.parse(v as String),
      ),
      updatedAt: $checkedConvert(
        'updated_at',
        (v) => v == null ? null : DateTime.parse(v as String),
      ),
      title: $checkedConvert('title', (v) => v as String?),
      desc: $checkedConvert('desc', (v) => v as String?),
      content: $checkedConvert('content', (v) => v as String?),
      img: $checkedConvert('img', (v) => v as String?),
      type: $checkedConvert('type', (v) => (v as num?)?.toInt()),
      status: $checkedConvert('status', (v) => (v as num?)?.toInt()),
      isTop: $checkedConvert('is_top', (v) => v as bool?),
      isDelete: $checkedConvert('is_delete', (v) => v as bool?),
      originalUrl: $checkedConvert('original_url', (v) => v as String?),
      categoryId: $checkedConvert('category_id', (v) => (v as num?)?.toInt()),
      tags: $checkedConvert(
        'tags',
        (v) => (v as List<dynamic>?)
            ?.map((e) => Category.fromJson(e as Map<String, dynamic>))
            .toList(),
      ),
      category: $checkedConvert(
        'category',
        (v) => v == null ? null : Category.fromJson(v as Map<String, dynamic>),
      ),
      user: $checkedConvert('user', (v) => v),
    );
    return val;
  },
  fieldKeyMap: const {
    'createdAt': 'created_at',
    'updatedAt': 'updated_at',
    'isTop': 'is_top',
    'isDelete': 'is_delete',
    'originalUrl': 'original_url',
    'categoryId': 'category_id',
  },
);

Map<String, dynamic> _$PostToJson(Post instance) => <String, dynamic>{
  'id': instance.id,
  'created_at': instance.createdAt?.toIso8601String(),
  'updated_at': instance.updatedAt?.toIso8601String(),
  'title': instance.title,
  'desc': instance.desc,
  'content': instance.content,
  'img': instance.img,
  'type': instance.type,
  'status': instance.status,
  'is_top': instance.isTop,
  'is_delete': instance.isDelete,
  'original_url': instance.originalUrl,
  'category_id': instance.categoryId,
  'tags': instance.tags,
  'category': instance.category,
  'user': instance.user,
};

Category _$CategoryFromJson(Map<String, dynamic> json) => $checkedCreate(
  'Category',
  json,
  ($checkedConvert) {
    final val = Category(
      id: $checkedConvert('id', (v) => (v as num).toInt()),
      createdAt: $checkedConvert(
        'created_at',
        (v) => DateTime.parse(v as String),
      ),
      updatedAt: $checkedConvert(
        'updated_at',
        (v) => DateTime.parse(v as String),
      ),
      name: $checkedConvert('name', (v) => v as String),
    );
    return val;
  },
  fieldKeyMap: const {'createdAt': 'created_at', 'updatedAt': 'updated_at'},
);

Map<String, dynamic> _$CategoryToJson(Category instance) => <String, dynamic>{
  'id': instance.id,
  'created_at': instance.createdAt.toIso8601String(),
  'updated_at': instance.updatedAt.toIso8601String(),
  'name': instance.name,
};
