/**************************************/
/*                                    */
/*   Visual Effects Engine - C++      */
/*     Frutiger Aero + Y2K Edition    */
/*           Programmed by            */
/*            Sertaç Ataç             */
/*            02.01.2026              */
/*                                    */
/**************************************/

#include <math.h>
#include <stdlib.h>
#include <string.h>
#include "ui_components.h"

/**************************************************/
/*                                                */
/*          GLOW EFFECT CALCULATIONS              */
/*                                                */
/**************************************************/

typedef struct {
    float intensity;
    float targetIntensity;
    float pulseSpeed;
    float time;
} GlowEffect;

GlowEffect create_glow_effect(float baseIntensity, float pulseSpeed) {
    GlowEffect effect;
    effect.intensity = baseIntensity;
    effect.targetIntensity = baseIntensity;
    effect.pulseSpeed = pulseSpeed;
    effect.time = 0.0f;
    return effect;
}

void update_glow_effect(GlowEffect* effect, float deltaTime) {
    effect->time += deltaTime * effect->pulseSpeed;
    float pulse = (sinf(effect->time) + 1.0f) * 0.5f;
    effect->intensity = effect->targetIntensity * (0.7f + pulse * 0.3f);
}

/**************************************************/
/*                                                */
/*             SCANLINE EFFECT                    */
/*                                                */
/**************************************************/

typedef struct {
    float position;
    float speed;
    float width;
    float alpha;
    int screenHeight;
} ScanlineEffect;

ScanlineEffect create_scanline_effect(int screenHeight) {
    ScanlineEffect effect;
    effect.position = 0.0f;
    effect.speed = 150.0f;
    effect.width = 2.0f;
    effect.alpha = 0.1f;
    effect.screenHeight = screenHeight;
    return effect;
}

/**************************************************/
/*                                                */
/*            NEON BORDER EFFECT                  */
/*                                                */
/**************************************************/

typedef struct {
    RGBAColor innerColor;
    RGBAColor outerColor;
    float thickness;
    float glowRadius;
    float pulseTime;
} NeonBorder;

NeonBorder create_neon_border(RGBAColor color, float thickness) {
    NeonBorder border;
    border.innerColor = color;
    border.outerColor = (RGBAColor){color.r, color.g, color.b, 60};
    border.thickness = thickness;
    border.glowRadius = thickness * 3.0f;
    border.pulseTime = 0.0f;
    return border;
}

/**************************************************/
/*                                                */
/*         EXPORTED CGO FUNCTIONS                 */
/*                                                */
/**************************************************/

#ifdef __cplusplus
extern "C" {
#endif

float GetGlowIntensity(float time, float baseIntensity) {
    float pulse = (sinf(time * 2.5f) + 1.0f) * 0.5f;
    return baseIntensity * (0.7f + pulse * 0.3f);
}

float GetScanlineAlpha(float y, float scanlineY, float width) {
    float dist = fabsf(y - scanlineY);
    if (dist < width) {
        return (1.0f - dist / width) * 0.15f;
    }
    return 0.0f;
}

float GetNeonPulse(float time) {
    return 0.6f + 0.4f * sinf(time * 3.0f);
}

#ifdef __cplusplus
}
#endif
