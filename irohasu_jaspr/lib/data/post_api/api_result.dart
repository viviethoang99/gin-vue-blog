import 'package:json_annotation/json_annotation.dart';

part 'api_result.g.dart';

@JsonSerializable(genericArgumentFactories: true, checked: true)
class ApiResult<T> {
  final int code;
  final String message;
  final T data;

  const ApiResult({required this.data, this.code = 0, this.message = ''});

  factory ApiResult.fromJson(Map<String, dynamic> json, T Function(Object?) fromJsonT) =>
      _$ApiResultFromJson(json, fromJsonT);

  Map<String, dynamic> toJson(Object? Function(T value) toJsonT) => _$ApiResultToJson(this, toJsonT);

  @override
  String toString() => 'ApiResult(code: $code, message: $message, data: $data)';
}
