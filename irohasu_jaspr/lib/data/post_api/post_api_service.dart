import 'package:dio/dio.dart';
import 'package:pretty_dio_logger/pretty_dio_logger.dart';
import 'package:retrofit/retrofit.dart';

import '../model/post.dart';
import 'api_result.dart';
import 'home_data.dart';
import 'list_posts.dart';

part 'post_api_service.g.dart';

@RestApi()
abstract class ApiService {
  factory ApiService(Dio dio, {String? baseUrl}) = _ApiService;

  factory ApiService.create({Dio? dio, String? baseUrl, bool enableLogging = false}) {
    final client =
        dio ??
        Dio(
          BaseOptions(
            connectTimeout: const Duration(seconds: 10),
            receiveTimeout: const Duration(seconds: 15),
            headers: const {'accept': 'application/json'},
          ),
        );
    if (enableLogging) {
      client.interceptors.add(
        PrettyDioLogger(
          requestHeader: true,
          requestBody: true,
          responseBody: true,
          responseHeader: false,
          error: true,
          compact: true,
          maxWidth: 90,
        ),
      );
    }
    return ApiService(client, baseUrl: baseUrl);
  }

  @POST('/api/login')
  Future<dynamic> login(@Body() Map<String, dynamic> body);

  @POST('/api/register')
  Future<dynamic> register(@Body() Map<String, dynamic> body);

  @GET('/api/logout')
  Future<dynamic> logout();

  @GET('/api/code')
  Future<dynamic> getCode();

  @GET('/api/front/about')
  Future<dynamic> getAbout();

  @GET('/api/front/page')
  Future<dynamic> getPage(@Queries() Map<String, dynamic>? queries);

  @GET('/api/front/home')
  Future<ApiResult<HomeData>> getHome();

  @GET('/api/front/article/list')
  Future<ListPosts> getPosts(@Queries() Map<String, dynamic>? queries);

  @GET('/api/front/article/{id}')
  Future<ApiResult<Post>> getContent(@Path('id') String id);

  @GET('/api/front/article/archive')
  Future<dynamic> getArticleArchive(@Queries() Map<String, dynamic>? queries);

  @GET('/api/front/article/search')
  Future<ListPosts> searchArticles(@Queries() Map<String, dynamic>? queries);

  @GET('/api/front/category/list')
  Future<dynamic> getCategoryList(@Queries() Map<String, dynamic>? queries);

  @GET('/api/front/tag/list')
  Future<dynamic> getTagList(@Queries() Map<String, dynamic>? queries);

  @GET('/api/front/message/list')
  Future<dynamic> getMessageList(@Queries() Map<String, dynamic>? queries);

  @GET('/api/front/link/list')
  Future<dynamic> getLinkList(@Queries() Map<String, dynamic>? queries);

  @GET('/api/front/comment/list')
  Future<dynamic> getCommentList(@Queries() Map<String, dynamic>? queries);

  @GET('/api/front/comment/replies/{id}')
  Future<dynamic> getCommentReplies(@Path('id') String id);

  @GET('/api/front/user/info')
  Future<dynamic> getUserInfo();

  @PUT('/api/front/user/info')
  Future<dynamic> updateUserInfo(@Body() Map<String, dynamic> body);

  @POST('/api/front/message')
  Future<dynamic> createMessage(@Body() Map<String, dynamic> body);

  @POST('/api/front/comment')
  Future<dynamic> createComment(@Body() Map<String, dynamic> body);

  @GET('/api/front/comment/like/{id}')
  Future<dynamic> likeComment(@Path('id') String id);

  @GET('/api/front/article/like/{id}')
  Future<dynamic> likeArticle(@Path('id') String id);
}
