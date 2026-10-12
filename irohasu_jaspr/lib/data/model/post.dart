import 'package:json_annotation/json_annotation.dart';

part 'post.g.dart';

typedef Posts = List<Post>;

@JsonSerializable(fieldRename: FieldRename.snake, checked: true)
final class Post {
  final int? id;
  final DateTime? createdAt;
  final DateTime? updatedAt;
  final String? title;
  final String? desc;
  final String? content;
  final String? img;
  final int? type;
  final int? status;
  final bool? isTop;
  final bool? isDelete;
  final String? originalUrl;
  final int? categoryId;
  final List<Category>? tags;
  final Category? category;
  final dynamic user;

  const Post({
    this.id,
    this.createdAt,
    this.updatedAt,
    this.title,
    this.desc,
    this.content,
    this.img,
    this.type,
    this.status,
    this.isTop,
    this.isDelete,
    this.originalUrl,
    this.categoryId,
    this.tags,
    this.category,
    this.user,
  });

  factory Post.fromJson(Map<String, dynamic> json) => _$PostFromJson(json);

  Map<String, dynamic> toJson() => _$PostToJson(this);
}

@JsonSerializable(fieldRename: FieldRename.snake, checked: true)
final class Category {
  const Category({
    required this.id,
    required this.createdAt,
    required this.updatedAt,
    required this.name,
  });

  final int id;
  final DateTime createdAt;
  final DateTime updatedAt;
  final String name;

  factory Category.fromJson(Map<String, dynamic> json) => _$CategoryFromJson(json);

  Map<String, dynamic> toJson() => _$CategoryToJson(this);
}
