import 'package:json_annotation/json_annotation.dart';

import '../model/post.dart';

part 'list_posts.g.dart';

@JsonSerializable(
  fieldRename: FieldRename.snake,
  checked: true,
  createFactory: false,
  createToJson: false,
)
class ListPosts {
  final int code;
  final String message;
  final List<Post>? data;
  final Paging? paging;
  final Filter? filter;

  const ListPosts({
    this.code = 0,
    this.message = '',
    this.data,
    this.paging,
    this.filter,
  });

  factory ListPosts.fromJson(Map<String, dynamic> json) {
    final rawData = json['data'];
    final pageData = rawData is Map ? rawData : null;

    return ListPosts(
      code: _asInt(json['code']) ?? 0,
      message: json['message'] as String? ?? '',
      data: _decodePosts(pageData?['page_data'] ?? rawData),
      paging: switch (json['paging']) {
        final Map<dynamic, dynamic> value => Paging.fromJson(Map<String, dynamic>.from(value)),
        _ when pageData != null => Paging(
          page: _asInt(pageData['page_num']),
          limit: _asInt(pageData['page_size']),
          total: _asInt(pageData['total']),
        ),
        _ => null,
      },
      filter: switch (json['filter']) {
        final Map<dynamic, dynamic> value => Filter.fromJson(Map<String, dynamic>.from(value)),
        _ => null,
      },
    );
  }

  Map<String, dynamic> toJson() => {
    'code': code,
    'message': message,
    'data': data?.map((post) => post.toJson()).toList(growable: false),
    'paging': paging?.toJson(),
    'filter': filter?.toJson(),
  };
}

List<Post>? _decodePosts(Object? value) {
  if (value == null) return null;
  if (value is! List) {
    throw FormatException('Expected article data to be a list, got ${value.runtimeType}.');
  }

  return value.map((item) => Post.fromJson(Map<String, dynamic>.from(item as Map))).toList(growable: false);
}

int? _asInt(Object? value) => value is num ? value.toInt() : int.tryParse('$value');

@JsonSerializable(fieldRename: FieldRename.snake, checked: true)
class Filter {
  final int? ownerId;

  const Filter({
    this.ownerId,
  });

  factory Filter.fromJson(Map<String, dynamic> json) => _$FilterFromJson(json);

  Map<String, dynamic> toJson() => _$FilterToJson(this);
}

@JsonSerializable(fieldRename: FieldRename.snake, checked: true)
class Paging {
  final int? page;
  final int? limit;
  final int? total;
  final String? cursor;
  final String? nextCursor;

  const Paging({
    this.page,
    this.limit,
    this.total,
    this.cursor,
    this.nextCursor,
  });

  factory Paging.fromJson(Map<String, dynamic> json) => _$PagingFromJson(json);

  Map<String, dynamic> toJson() => _$PagingToJson(this);
}
