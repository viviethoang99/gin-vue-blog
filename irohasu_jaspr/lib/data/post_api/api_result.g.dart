// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'api_result.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ApiResult<T> _$ApiResultFromJson<T>(
  Map<String, dynamic> json,
  T Function(Object? json) fromJsonT,
) => $checkedCreate('ApiResult', json, ($checkedConvert) {
  final val = ApiResult<T>(
    data: $checkedConvert('data', (v) => fromJsonT(v)),
    code: $checkedConvert('code', (v) => (v as num?)?.toInt() ?? 0),
    message: $checkedConvert('message', (v) => v as String? ?? ''),
  );
  return val;
});

Map<String, dynamic> _$ApiResultToJson<T>(
  ApiResult<T> instance,
  Object? Function(T value) toJsonT,
) => <String, dynamic>{
  'code': instance.code,
  'message': instance.message,
  'data': toJsonT(instance.data),
};
