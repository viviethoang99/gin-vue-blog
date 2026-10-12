import 'dart:async';
import 'dart:math';

import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';
import 'package:universal_web/js_interop.dart';
import 'package:universal_web/web.dart' as web;

/// Configuration used to build the static plexus network.
class StaticParticleConfig {
  const StaticParticleConfig({
    this.nodeCount = 110,
    this.maxNodeCount = 170,
    this.nodeCountReferenceArea = 2073600,
    this.connectionProbability = 0.54,
    this.connectionDistanceScale = 0.17,
    this.minConnectionDistance = 130,
    this.maxConnectionDistance = 280,
    this.viewportOverscan = 0.08,
    this.lineWidth = 0.5,
    this.lightColor = '0, 0, 0',
    this.darkColor = '255, 255, 255',
    this.lightMinOpacity = 0.06,
    this.lightMaxOpacity = 0.12,
    this.darkMinOpacity = 0.07,
    this.darkMaxOpacity = 0.15,
    this.animationDurationMs = 750,
    this.resizeDebounceMs = 180,
  });

  final int nodeCount;
  final int maxNodeCount;
  final double nodeCountReferenceArea;
  final double connectionProbability;
  final double connectionDistanceScale;
  final double minConnectionDistance;
  final double maxConnectionDistance;
  final double viewportOverscan;
  final double lineWidth;
  final String lightColor;
  final String darkColor;
  final double lightMinOpacity;
  final double lightMaxOpacity;
  final double darkMinOpacity;
  final double darkMaxOpacity;
  final int animationDurationMs;
  final int resizeDebounceMs;

  int effectiveNodeCount(double width, double height) {
    if (nodeCount <= 0 || width <= 0 || height <= 0 || nodeCountReferenceArea <= 0) return 0;

    final areaScale = sqrt((width * height) / nodeCountReferenceArea);
    final scaledNodeCount = (nodeCount * max(1.0, areaScale)).round();
    return scaledNodeCount.clamp(nodeCount, max(maxNodeCount, nodeCount));
  }

  double connectionDistance(double width, double height) {
    final diagonal = sqrt(width * width + height * height);
    return (diagonal * connectionDistanceScale).clamp(minConnectionDistance, maxConnectionDistance);
  }
}

class _ParticlePoint {
  const _ParticlePoint(this.x, this.y);

  final double x;
  final double y;
}

class _ParticleEdge {
  const _ParticleEdge(this.start, this.end, this.opacityFactor);

  final _ParticlePoint start;
  final _ParticlePoint end;
  final double opacityFactor;
}

/// A non-interactive, viewport-sized canvas containing a static plexus network.
///
/// Every constructor parameter is primitive so Jaspr can serialize this
/// client component safely during static pre-rendering.
@client
class StaticParticleBackground extends StatefulComponent {
  const StaticParticleBackground({
    this.nodeCount = 110,
    this.maxNodeCount = 170,
    this.nodeCountReferenceArea = 2073600,
    this.connectionProbability = 0.54,
    this.connectionDistanceScale = 0.17,
    this.minConnectionDistance = 130,
    this.maxConnectionDistance = 280,
    this.viewportOverscan = 0.08,
    this.lineWidth = 0.5,
    this.lightColor = '0, 0, 0',
    this.darkColor = '255, 255, 255',
    this.lightMinOpacity = 0.06,
    this.lightMaxOpacity = 0.12,
    this.darkMinOpacity = 0.07,
    this.darkMaxOpacity = 0.15,
    this.animationDurationMs = 750,
    this.resizeDebounceMs = 180,
    super.key,
  });

  final int nodeCount;
  final int maxNodeCount;
  final double nodeCountReferenceArea;
  final double connectionProbability;
  final double connectionDistanceScale;
  final double minConnectionDistance;
  final double maxConnectionDistance;
  final double viewportOverscan;
  final double lineWidth;
  final String lightColor;
  final String darkColor;
  final double lightMinOpacity;
  final double lightMaxOpacity;
  final double darkMinOpacity;
  final double darkMaxOpacity;
  final int animationDurationMs;
  final int resizeDebounceMs;

  StaticParticleConfig get config => StaticParticleConfig(
    nodeCount: nodeCount,
    maxNodeCount: maxNodeCount,
    nodeCountReferenceArea: nodeCountReferenceArea,
    connectionProbability: connectionProbability,
    connectionDistanceScale: connectionDistanceScale,
    minConnectionDistance: minConnectionDistance,
    maxConnectionDistance: maxConnectionDistance,
    viewportOverscan: viewportOverscan,
    lineWidth: lineWidth,
    lightColor: lightColor,
    darkColor: darkColor,
    lightMinOpacity: lightMinOpacity,
    lightMaxOpacity: lightMaxOpacity,
    darkMinOpacity: darkMinOpacity,
    darkMaxOpacity: darkMaxOpacity,
    animationDurationMs: animationDurationMs,
    resizeDebounceMs: resizeDebounceMs,
  );

  @override
  State<StaticParticleBackground> createState() => _StaticParticleBackgroundState();

  @css
  static List<StyleRule> get styles => [
    css('.static-particle-background').styles(
      display: .block,
      position: .fixed(top: 0.px, right: 0.px, bottom: 0.px, left: 0.px),
      zIndex: ZIndex(0),
      width: 100.vw,
      height: 100.vh,
      pointerEvents: .none,
    ),
  ];
}

class _StaticParticleBackgroundState extends State<StaticParticleBackground> {
  final GlobalNodeKey<web.HTMLCanvasElement> _canvasKey = GlobalNodeKey(debugLabel: 'plexus-canvas');
  final Random _random = Random();

  StreamSubscription<web.Event>? _resizeSubscription;
  StreamSubscription<web.Event>? _colorSchemeSubscription;
  Timer? _setupTimer;
  Timer? _resizeTimer;
  web.MediaQueryList? _colorSchemeQuery;
  web.MutationObserver? _themeObserver;
  List<_ParticleEdge> _edges = const [];
  int? _animationFrameId;
  double? _animationStart;
  double _width = 0;
  double _height = 0;
  bool _disposed = false;

  StaticParticleConfig get _config => component.config;

  @override
  void initState() {
    super.initState();
    if (!kIsWeb) return;

    // Queue setup so the GlobalNodeKey has its hydrated DOM node.
    _setupTimer = Timer(Duration.zero, _initializeCanvas);
  }

  void _initializeCanvas() {
    if (_disposed) return;

    _resizeSubscription = web.EventStreamProviders.resizeEvent.forTarget(web.window).listen((_) {
      _resizeTimer?.cancel();
      _resizeTimer = Timer(Duration(milliseconds: _config.resizeDebounceMs), () => _renderNetwork(animate: false));
    });

    _colorSchemeQuery = web.window.matchMedia('(prefers-color-scheme: dark)');
    _colorSchemeSubscription = web.EventStreamProviders.changeEvent.forTarget(_colorSchemeQuery).listen((_) {
      _redrawForTheme();
    });

    final documentRoot = web.document.documentElement;
    if (documentRoot != null) {
      _themeObserver = web.MutationObserver(
        ((JSArray<web.MutationRecord> _, web.MutationObserver observer) => _redrawForTheme()).toJS,
      )..observe(documentRoot, web.MutationObserverInit(attributes: true));
    }

    _renderNetwork(animate: !_prefersReducedMotion());
  }

  void _renderNetwork({required bool animate}) {
    if (_disposed) return;

    final canvas = _canvasKey.currentNode;
    if (canvas == null) return;

    _cancelAnimation();
    _width = web.window.innerWidth.toDouble();
    _height = web.window.innerHeight.toDouble();
    if (_width <= 0 || _height <= 0) return;

    final pixelRatio = web.window.devicePixelRatio.clamp(1.0, 3.0);
    canvas.width = (_width * pixelRatio).round();
    canvas.height = (_height * pixelRatio).round();

    final context = canvas.context2D;
    context.resetTransform();
    context.scale(pixelRatio, pixelRatio);

    final points = _generatePoints(_width, _height);
    _edges = _generateEdges(points, _config.connectionDistance(_width, _height));

    if (!animate || _config.animationDurationMs <= 0) {
      _draw(progress: 1);
      return;
    }

    _animationStart = null;
    _animationFrameId = web.window.requestAnimationFrame(_animate.toJS);
  }

  List<_ParticlePoint> _generatePoints(double width, double height) {
    final nodeCount = _config.effectiveNodeCount(width, height);
    if (nodeCount <= 0) return const [];

    final points = <_ParticlePoint>[];
    final overscanX = width * _config.viewportOverscan;
    final overscanY = height * _config.viewportOverscan;
    final outsideCount = min(nodeCount, max(4, nodeCount ~/ 8));

    for (var index = 0; index < outsideCount; index++) {
      switch (index % 4) {
        case 0:
          points.add(_ParticlePoint(-_random.nextDouble() * overscanX, _random.nextDouble() * height));
        case 1:
          points.add(_ParticlePoint(width + _random.nextDouble() * overscanX, _random.nextDouble() * height));
        case 2:
          points.add(_ParticlePoint(_random.nextDouble() * width, -_random.nextDouble() * overscanY));
        case 3:
          points.add(_ParticlePoint(_random.nextDouble() * width, height + _random.nextDouble() * overscanY));
      }
    }

    while (points.length < nodeCount) {
      points.add(_ParticlePoint(_random.nextDouble() * width, _random.nextDouble() * height));
    }
    return points;
  }

  List<_ParticleEdge> _generateEdges(List<_ParticlePoint> points, double maxDistance) {
    final edges = <_ParticleEdge>[];
    final maxDistanceSquared = maxDistance * maxDistance;

    for (var first = 0; first < points.length; first++) {
      for (var second = first + 1; second < points.length; second++) {
        final dx = points[first].x - points[second].x;
        final dy = points[first].y - points[second].y;
        if (dx * dx + dy * dy > maxDistanceSquared || _random.nextDouble() > _config.connectionProbability) {
          continue;
        }

        edges.add(_ParticleEdge(points[first], points[second], _random.nextDouble()));
      }
    }
    return edges;
  }

  void _animate(double timestamp) {
    if (_disposed) return;

    _animationStart ??= timestamp;
    final elapsed = timestamp - _animationStart!;
    final linearProgress = (elapsed / _config.animationDurationMs).clamp(0.0, 1.0);
    final easedProgress = 1 - pow(1 - linearProgress, 3).toDouble();
    _draw(progress: easedProgress);

    if (linearProgress < 1) {
      _animationFrameId = web.window.requestAnimationFrame(_animate.toJS);
    } else {
      _animationFrameId = null;
      _animationStart = null;
    }
  }

  void _draw({required double progress}) {
    if (_disposed || _width <= 0 || _height <= 0) return;

    final canvas = _canvasKey.currentNode;
    if (canvas == null) return;
    final context = canvas.context2D;
    context.clearRect(0, 0, _width, _height);
    context.lineWidth = _config.lineWidth;

    final dark = _isDarkTheme();
    final color = dark ? _config.darkColor : _config.lightColor;
    final minOpacity = dark ? _config.darkMinOpacity : _config.lightMinOpacity;
    final maxOpacity = dark ? _config.darkMaxOpacity : _config.lightMaxOpacity;

    for (final edge in _edges) {
      final opacity = (minOpacity + (maxOpacity - minOpacity) * edge.opacityFactor) * progress;
      context.strokeStyle = 'rgba($color, $opacity)'.toJS;
      context.beginPath();
      context.moveTo(edge.start.x, edge.start.y);
      context.lineTo(edge.end.x, edge.end.y);
      context.stroke();
    }
  }

  void _redrawForTheme() {
    _cancelAnimation();
    _draw(progress: 1);
  }

  bool _isDarkTheme() {
    final root = web.document.documentElement;
    final declaredTheme = root?.getAttribute('data-theme')?.toLowerCase();
    if (declaredTheme == 'dark' || root?.classList.contains('dark') == true) return true;
    if (declaredTheme == 'light' || root?.classList.contains('light') == true) return false;

    // A system preference is not necessarily the theme used by the page. The
    // current Irohasu layout, for example, is light even when the OS is dark.
    // Reading the painted background keeps the lines contrasting correctly and
    // still responds to CSS media-query theme changes.
    return _elementHasDarkBackground(web.document.body) || _elementHasDarkBackground(root);
  }

  bool _elementHasDarkBackground(web.Element? element) {
    if (element == null) return false;

    final color = web.window.getComputedStyle(element).backgroundColor;
    final channels = RegExp(r'[\d.]+').allMatches(color).map((match) => double.parse(match.group(0)!)).toList();
    if (channels.length < 3 || (channels.length > 3 && channels[3] == 0)) return false;

    final luminance = (channels[0] * 0.2126 + channels[1] * 0.7152 + channels[2] * 0.0722) / 255;
    return luminance < 0.5;
  }

  bool _prefersReducedMotion() => web.window.matchMedia('(prefers-reduced-motion: reduce)').matches;

  void _cancelAnimation() {
    final animationFrameId = _animationFrameId;
    if (animationFrameId != null) {
      web.window.cancelAnimationFrame(animationFrameId);
      _animationFrameId = null;
    }
    _animationStart = null;
  }

  @override
  void dispose() {
    _disposed = true;
    _cancelAnimation();
    _setupTimer?.cancel();
    _resizeTimer?.cancel();
    _resizeSubscription?.cancel();
    _colorSchemeSubscription?.cancel();
    _themeObserver?.disconnect();
    super.dispose();
  }

  @override
  Component build(BuildContext context) {
    return Component.element(
      tag: 'canvas',
      key: _canvasKey,
      classes: 'static-particle-background',
      attributes: const {'aria-hidden': 'true'},
    );
  }
}
