import 'package:flutter_test/flutter_test.dart';
import 'package:irohasu_jaspr/components/static_particle_background.dart';

void main() {
  group('StaticParticleConfig.connectionDistance', () {
    const config = StaticParticleConfig();

    test('uses the minimum distance on small viewports', () {
      expect(config.connectionDistance(320, 568), config.minConnectionDistance);
    });

    test('scales with a medium viewport', () {
      final distance = config.connectionDistance(1024, 768);

      expect(distance, closeTo(217.6, 0.001));
    });

    test('uses the maximum distance on large viewports', () {
      expect(config.connectionDistance(2560, 1440), config.maxConnectionDistance);
    });
  });

  group('StaticParticleConfig.effectiveNodeCount', () {
    const config = StaticParticleConfig();

    test('keeps the base count on mobile and full HD viewports', () {
      expect(config.effectiveNodeCount(390, 844), 110);
      expect(config.effectiveNodeCount(1920, 1080), 110);
    });

    test('increases the count for a 2K viewport', () {
      expect(config.effectiveNodeCount(2560, 1440), 147);
    });

    test('caps the count on very large viewports', () {
      expect(config.effectiveNodeCount(3840, 2160), 170);
    });
  });
}
