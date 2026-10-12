// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'list_posts.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Filter _$FilterFromJson(Map<String, dynamic> json) =>
    $checkedCreate('Filter', json, ($checkedConvert) {
      final val = Filter(
        ownerId: $checkedConvert('owner_id', (v) => (v as num?)?.toInt()),
      );
      return val;
    }, fieldKeyMap: const {'ownerId': 'owner_id'});

Map<String, dynamic> _$FilterToJson(Filter instance) => <String, dynamic>{
  'owner_id': instance.ownerId,
};

Paging _$PagingFromJson(Map<String, dynamic> json) =>
    $checkedCreate('Paging', json, ($checkedConvert) {
      final val = Paging(
        page: $checkedConvert('page', (v) => (v as num?)?.toInt()),
        limit: $checkedConvert('limit', (v) => (v as num?)?.toInt()),
        total: $checkedConvert('total', (v) => (v as num?)?.toInt()),
        cursor: $checkedConvert('cursor', (v) => v as String?),
        nextCursor: $checkedConvert('next_cursor', (v) => v as String?),
      );
      return val;
    }, fieldKeyMap: const {'nextCursor': 'next_cursor'});

Map<String, dynamic> _$PagingToJson(Paging instance) => <String, dynamic>{
  'page': instance.page,
  'limit': instance.limit,
  'total': instance.total,
  'cursor': instance.cursor,
  'next_cursor': instance.nextCursor,
};
