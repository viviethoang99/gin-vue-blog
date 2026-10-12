import '../../config/environment.dart';
import '../model/post.dart';
import 'api_result.dart';
import 'home_data.dart';
import 'list_posts.dart';
import 'post_api_service.dart';

class PostRepository {
  PostRepository(this._api);

  factory PostRepository.live({String? baseUrl, bool enableLogging = false}) => PostRepository(
    ApiService.create(
      baseUrl: baseUrl ?? apiBaseUrl,
      enableLogging: enableLogging,
    ),
  );

  final ApiService _api;

  Future<ApiResult<HomeData>> getHome() => _api.getHome();

  Future<ListPosts> getPosts({int? page, int? limit}) {
    return _api.getPosts({
      'page_num': ?page,
      'page_size': ?limit,
    });
  }

  Future<ApiResult<Post>> getPost(int id) => _api.getContent('$id');

  Future<ListPosts> searchPosts(String keyword) {
    return _api.searchArticles({'keyword': keyword});
  }
}
