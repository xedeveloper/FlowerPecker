package components

import (
	"os"
	"path/filepath"
)

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func GenerateUIComponents(projectPath string) error {
	widgetsBase := filepath.Join(projectPath, "lib/core/widgets")

	files := map[string]string{
		filepath.Join(widgetsBase, "fp_button.dart"):     fpButtonDart(),
		filepath.Join(widgetsBase, "fp_input.dart"):      fpInputDart(),
		filepath.Join(widgetsBase, "fp_card.dart"):       fpCardDart(),
		filepath.Join(widgetsBase, "fp_nav.dart"):        fpNavDart(),
		filepath.Join(widgetsBase, "fp_footer.dart"):     fpFooterDart(),
		filepath.Join(widgetsBase, "fp_hero.dart"):       fpHeroDart(),
		filepath.Join(widgetsBase, "fp_divider.dart"):    fpDividerDart(),
		filepath.Join(widgetsBase, "fp_toast.dart"):      fpToastDart(),
		filepath.Join(widgetsBase, "fp_typography.dart"): fpTypographyDart(),
		filepath.Join(widgetsBase, "widgets.dart"):       widgetsBarrelDart(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

func widgetsBarrelDart() string {
	return `export 'fp_button.dart';
export 'fp_input.dart';
export 'fp_card.dart';
export 'fp_nav.dart';
export 'fp_footer.dart';
export 'fp_hero.dart';
export 'fp_divider.dart';
export 'fp_toast.dart';
export 'fp_typography.dart';
`
}

func fpButtonDart() string {
	return `import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

class FPButtonPrimary extends StatefulWidget {
  final String label;
  final VoidCallback? onPressed;
  final bool isLoading;

  const FPButtonPrimary({
    super.key,
    required this.label,
    this.onPressed,
    this.isLoading = false,
  });

  @override
  State<FPButtonPrimary> createState() => _FPButtonPrimaryState();
}

class _FPButtonPrimaryState extends State<FPButtonPrimary> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: GestureDetector(
        onTap: widget.isLoading ? null : widget.onPressed,
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 150),
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
          decoration: BoxDecoration(
            color: _isHovered ? AppColors.inkSoft : AppColors.primary,
          ),
          child: widget.isLoading
              ? const SizedBox(
                  width: 16,
                  height: 16,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: AppColors.onPrimary,
                  ),
                )
              : Text(
                  widget.label,
                  style: AppTextStyles.labelMd.copyWith(
                    color: AppColors.onPrimary,
                  ),
                ),
        ),
      ),
    );
  }
}

class FPButtonOutline extends StatefulWidget {
  final String label;
  final VoidCallback? onPressed;

  const FPButtonOutline({
    super.key,
    required this.label,
    this.onPressed,
  });

  @override
  State<FPButtonOutline> createState() => _FPButtonOutlineState();
}

class _FPButtonOutlineState extends State<FPButtonOutline> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: GestureDetector(
        onTap: widget.onPressed,
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 150),
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
          decoration: BoxDecoration(
            color: _isHovered ? AppColors.canvasSoft : AppColors.canvas,
            border: Border.all(color: AppColors.primary),
          ),
          child: Text(
            widget.label,
            style: AppTextStyles.labelMd.copyWith(color: AppColors.ink),
          ),
        ),
      ),
    );
  }
}

class FPButtonIconCircular extends StatelessWidget {
  final IconData icon;
  final VoidCallback? onPressed;
  final double size;

  const FPButtonIconCircular({
    super.key,
    required this.icon,
    this.onPressed,
    this.size = 48,
  });

  @override
  Widget build(BuildContext context) {
    return Material(
      color: AppColors.canvas,
      shape: const CircleBorder(
        side: BorderSide(color: AppColors.hairline),
      ),
      child: InkWell(
        onTap: onPressed,
        customBorder: const CircleBorder(),
        child: SizedBox(
          width: size,
          height: size,
          child: Icon(icon, color: AppColors.ink, size: size * 0.45),
        ),
      ),
    );
  }
}
`
}

func fpInputDart() string {
	return `import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

class FPTextInput extends StatefulWidget {
  final String label;
  final String? hint;
  final TextEditingController? controller;
  final ValueChanged<String>? onChanged;
  final String? errorText;
  final bool obscureText;
  final TextInputType keyboardType;
  final int? maxLines;

  const FPTextInput({
    super.key,
    required this.label,
    this.hint,
    this.controller,
    this.onChanged,
    this.errorText,
    this.obscureText = false,
    this.keyboardType = TextInputType.text,
    this.maxLines = 1,
  });

  @override
  State<FPTextInput> createState() => _FPTextInputState();
}

class _FPTextInputState extends State<FPTextInput> {
  bool _isFocused = false;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(widget.label, style: AppTextStyles.labelMd),
        const SizedBox(height: 6),
        Focus(
          onFocusChange: (focused) => setState(() => _isFocused = focused),
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 150),
            decoration: BoxDecoration(
              border: Border.all(
                color: _isFocused ? AppColors.ink : AppColors.hairline,
                width: _isFocused ? 2 : 1,
              ),
            ),
            child: TextField(
              controller: widget.controller,
              onChanged: widget.onChanged,
              obscureText: widget.obscureText,
              keyboardType: widget.keyboardType,
              maxLines: widget.maxLines,
              style: AppTextStyles.bodyMd,
              decoration: InputDecoration(
                hintText: widget.hint,
                hintStyle: AppTextStyles.bodyMd.copyWith(color: AppColors.body),
                border: InputBorder.none,
                contentPadding: const EdgeInsets.symmetric(
                  horizontal: 12,
                  vertical: 10,
                ),
              ),
            ),
          ),
        ),
        if (widget.errorText != null) ...[
          const SizedBox(height: 4),
          Text(
            widget.errorText!,
            style: AppTextStyles.bodySm.copyWith(color: Colors.red),
          ),
        ],
      ],
    );
  }
}
`
}

func fpCardDart() string {
	return `import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

class FPStoryCardLarge extends StatelessWidget {
  final String title;
  final String? subtitle;
  final String? imageUrl;
  final String? category;
  final VoidCallback? onTap;

  const FPStoryCardLarge({
    super.key,
    required this.title,
    this.subtitle,
    this.imageUrl,
    this.category,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        color: AppColors.canvas,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (imageUrl != null)
              AspectRatio(
                aspectRatio: 16 / 9,
                child: Image.network(
                  imageUrl!,
                  fit: BoxFit.cover,
                  errorBuilder: (_, __, ___) => Container(
                    color: AppColors.canvasSoft,
                    child: const Icon(Icons.image, color: AppColors.hairline),
                  ),
                ),
              ),
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (category != null) ...[
                    Text(category!.toUpperCase(), style: AppTextStyles.labelSm),
                    const SizedBox(height: 8),
                  ],
                  Text(title, style: AppTextStyles.displayMd),
                  if (subtitle != null) ...[
                    const SizedBox(height: 8),
                    Text(subtitle!, style: AppTextStyles.bodyMd),
                  ],
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class FPStoryCard extends StatelessWidget {
  final String title;
  final String? subtitle;
  final String? imageUrl;
  final VoidCallback? onTap;

  const FPStoryCard({
    super.key,
    required this.title,
    this.subtitle,
    this.imageUrl,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        color: AppColors.canvas,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (imageUrl != null)
              AspectRatio(
                aspectRatio: 4 / 3,
                child: Image.network(
                  imageUrl!,
                  fit: BoxFit.cover,
                  errorBuilder: (_, __, ___) => Container(
                    color: AppColors.canvasSoft,
                  ),
                ),
              ),
            Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: AppTextStyles.displayXs),
                  if (subtitle != null) ...[
                    const SizedBox(height: 4),
                    Text(subtitle!, style: AppTextStyles.bodySm),
                  ],
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class FPStoryRow extends StatefulWidget {
  final String title;
  final String? subtitle;
  final String? imageUrl;
  final VoidCallback? onTap;

  const FPStoryRow({
    super.key,
    required this.title,
    this.subtitle,
    this.imageUrl,
    this.onTap,
  });

  @override
  State<FPStoryRow> createState() => _FPStoryRowState();
}

class _FPStoryRowState extends State<FPStoryRow>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _opacity;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 400),
    );
    _opacity = Tween<double>(begin: 0, end: 1).animate(_controller);
    _controller.forward();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return FadeTransition(
      opacity: _opacity,
      child: GestureDetector(
        onTap: widget.onTap,
        child: Container(
          decoration: const BoxDecoration(
            border: Border(
              bottom: BorderSide(color: AppColors.hairline),
            ),
          ),
          padding: const EdgeInsets.symmetric(vertical: 16),
          child: Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(widget.title, style: AppTextStyles.displayXs),
                    if (widget.subtitle != null) ...[
                      const SizedBox(height: 4),
                      Text(widget.subtitle!, style: AppTextStyles.bodySm),
                    ],
                  ],
                ),
              ),
              if (widget.imageUrl != null)
                SizedBox(
                  width: 80,
                  height: 60,
                  child: Image.network(
                    widget.imageUrl!,
                    fit: BoxFit.cover,
                    errorBuilder: (_, __, ___) =>
                        Container(color: AppColors.canvasSoft),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
`
}

func fpNavDart() string {
	return `import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

class FPNavBar extends StatelessWidget implements PreferredSizeWidget {
  final String title;
  final List<FPNavLink> links;
  final Widget? trailing;

  const FPNavBar({
    super.key,
    required this.title,
    this.links = const [],
    this.trailing,
  });

  @override
  Size get preferredSize => const Size.fromHeight(56);

  @override
  Widget build(BuildContext context) {
    final isMobile = MediaQuery.of(context).size.width < 768;

    return Container(
      color: AppColors.canvas,
      child: Column(
        children: [
          Expanded(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: Row(
                children: [
                  Text(title.toUpperCase(), style: AppTextStyles.labelMd),
                  const Spacer(),
                  if (!isMobile) ...links,
                  if (isMobile)
                    IconButton(
                      icon: const Icon(Icons.menu, color: AppColors.ink),
                      onPressed: () => Scaffold.of(context).openEndDrawer(),
                    ),
                  if (trailing != null && !isMobile) trailing!,
                ],
              ),
            ),
          ),
          const Divider(height: 1, color: AppColors.hairline),
        ],
      ),
    );
  }
}

class FPNavLink extends StatefulWidget {
  final String label;
  final VoidCallback? onTap;
  final bool isActive;

  const FPNavLink({
    super.key,
    required this.label,
    this.onTap,
    this.isActive = false,
  });

  @override
  State<FPNavLink> createState() => _FPNavLinkState();
}

class _FPNavLinkState extends State<FPNavLink> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: GestureDetector(
        onTap: widget.onTap,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                widget.label,
                style: AppTextStyles.bodySmStrong.copyWith(
                  color: widget.isActive ? AppColors.ink : AppColors.body,
                ),
              ),
              AnimatedContainer(
                duration: const Duration(milliseconds: 150),
                height: 1,
                width: _isHovered || widget.isActive ? 40 : 0,
                color: AppColors.ink,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
`
}

func fpFooterDart() string {
	return `import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

class FPFooter extends StatelessWidget {
  final String brandName;
  final List<FPFooterColumn> columns;
  final String? copyright;

  const FPFooter({
    super.key,
    required this.brandName,
    this.columns = const [],
    this.copyright,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      color: AppColors.primary,
      padding: const EdgeInsets.symmetric(horizontal: 48, vertical: 48),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                flex: 2,
                child: Text(
                  brandName.toUpperCase(),
                  style: AppTextStyles.displaySm.copyWith(
                    color: AppColors.onPrimary,
                  ),
                ),
              ),
              ...columns.map(
                (col) => Expanded(child: col),
              ),
            ],
          ),
          const SizedBox(height: 48),
          const Divider(color: Color(0xFF333333)),
          const SizedBox(height: 16),
          Text(
            copyright ?? '© ${DateTime.now().year} $brandName. All rights reserved.',
            style: AppTextStyles.bodySm.copyWith(color: const Color(0xFF999999)),
          ),
        ],
      ),
    );
  }
}

class FPFooterColumn extends StatelessWidget {
  final String heading;
  final List<FPFooterLink> links;

  const FPFooterColumn({
    super.key,
    required this.heading,
    this.links = const [],
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          heading.toUpperCase(),
          style: AppTextStyles.labelSm.copyWith(color: const Color(0xFF999999)),
        ),
        const SizedBox(height: 16),
        ...links,
      ],
    );
  }
}

class FPFooterLink extends StatelessWidget {
  final String label;
  final VoidCallback? onTap;

  const FPFooterLink({super.key, required this.label, this.onTap});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: GestureDetector(
        onTap: onTap,
        child: Text(
          label,
          style: AppTextStyles.bodySm.copyWith(color: AppColors.onPrimary),
        ),
      ),
    );
  }
}
`
}

func fpHeroDart() string {
	return `import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

class FPHeroBand extends StatefulWidget {
  final String headline;
  final String? subheadline;
  final Widget? action;
  final Color backgroundColor;

  const FPHeroBand({
    super.key,
    required this.headline,
    this.subheadline,
    this.action,
    this.backgroundColor = AppColors.canvas,
  });

  @override
  State<FPHeroBand> createState() => _FPHeroBandState();
}

class _FPHeroBandState extends State<FPHeroBand>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _opacity;
  late Animation<Offset> _slide;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 800),
    );
    _opacity = Tween<double>(begin: 0, end: 1).animate(
      CurvedAnimation(parent: _controller, curve: Curves.easeOut),
    );
    _slide = Tween<Offset>(
      begin: const Offset(0, 0.1),
      end: Offset.zero,
    ).animate(CurvedAnimation(parent: _controller, curve: Curves.easeOut));
    _controller.forward();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      color: widget.backgroundColor,
      padding: const EdgeInsets.symmetric(horizontal: 48, vertical: 80),
      child: FadeTransition(
        opacity: _opacity,
        child: SlideTransition(
          position: _slide,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(widget.headline, style: AppTextStyles.displayHero),
              if (widget.subheadline != null) ...[
                const SizedBox(height: 24),
                Text(widget.subheadline!, style: AppTextStyles.bodyLg),
              ],
              if (widget.action != null) ...[
                const SizedBox(height: 40),
                widget.action!,
              ],
            ],
          ),
        ),
      ),
    );
  }
}
`
}

func fpDividerDart() string {
	return `import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

class FPHairlineDivider extends StatelessWidget {
  final double indent;
  final double endIndent;

  const FPHairlineDivider({
    super.key,
    this.indent = 0,
    this.endIndent = 0,
  });

  @override
  Widget build(BuildContext context) {
    return Divider(
      height: 1,
      thickness: 1,
      color: AppColors.hairline,
      indent: indent,
      endIndent: endIndent,
    );
  }
}
`
}

func fpToastDart() string {
	return `import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

enum FPToastType { info, success, error }

class FPToast extends StatefulWidget {
  final String message;
  final FPToastType type;
  final Duration duration;
  final VoidCallback? onDismissed;

  const FPToast({
    super.key,
    required this.message,
    this.type = FPToastType.info,
    this.duration = const Duration(seconds: 3),
    this.onDismissed,
  });

  static void show(
    BuildContext context,
    String message, {
    FPToastType type = FPToastType.info,
  }) {
    final overlay = Overlay.of(context);
    late OverlayEntry entry;
    entry = OverlayEntry(
      builder: (_) => Positioned(
        bottom: 24,
        left: 0,
        right: 0,
        child: Center(
          child: FPToast(
            message: message,
            type: type,
            onDismissed: () => entry.remove(),
          ),
        ),
      ),
    );
    overlay.insert(entry);
  }

  @override
  State<FPToast> createState() => _FPToastState();
}

class _FPToastState extends State<FPToast>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _opacity;
  late Animation<Offset> _slide;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 300),
    );
    _opacity = Tween<double>(begin: 0, end: 1).animate(_controller);
    _slide = Tween<Offset>(
      begin: const Offset(0, 0.5),
      end: Offset.zero,
    ).animate(CurvedAnimation(parent: _controller, curve: Curves.easeOut));

    _controller.forward();

    Future.delayed(widget.duration, () {
      if (mounted) {
        _controller.reverse().then((_) => widget.onDismissed?.call());
      }
    });
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Color get _backgroundColor {
    switch (widget.type) {
      case FPToastType.success:
        return const Color(0xFF2E7D32);
      case FPToastType.error:
        return const Color(0xFFD32F2F);
      case FPToastType.info:
        return AppColors.inkSoft;
    }
  }

  @override
  Widget build(BuildContext context) {
    return FadeTransition(
      opacity: _opacity,
      child: SlideTransition(
        position: _slide,
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
          decoration: BoxDecoration(
            color: _backgroundColor,
          ),
          child: Text(
            widget.message,
            style: AppTextStyles.bodySm.copyWith(color: AppColors.onPrimary),
          ),
        ),
      ),
    );
  }
}
`
}

func fpTypographyDart() string {
	return `import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

class FPText extends StatelessWidget {
  final String text;
  final TextStyle style;
  final TextAlign? textAlign;
  final int? maxLines;
  final TextOverflow? overflow;

  const FPText(
    this.text, {
    super.key,
    required this.style,
    this.textAlign,
    this.maxLines,
    this.overflow,
  });

  factory FPText.displayHero(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.displayHero, textAlign: textAlign);

  factory FPText.displayLg(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.displayLg, textAlign: textAlign);

  factory FPText.displayMd(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.displayMd, textAlign: textAlign);

  factory FPText.displaySm(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.displaySm, textAlign: textAlign);

  factory FPText.displayXs(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.displayXs, textAlign: textAlign);

  factory FPText.bodyLg(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.bodyLg, textAlign: textAlign);

  factory FPText.bodyMd(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.bodyMd, textAlign: textAlign);

  factory FPText.bodySm(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.bodySm, textAlign: textAlign);

  factory FPText.bodySmStrong(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.bodySmStrong, textAlign: textAlign);

  factory FPText.labelMd(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.labelMd, textAlign: textAlign);

  factory FPText.labelSm(String text, {Key? key, TextAlign? textAlign}) =>
      FPText(text, key: key, style: AppTextStyles.labelSm, textAlign: textAlign);

  @override
  Widget build(BuildContext context) {
    return Text(
      text,
      style: style,
      textAlign: textAlign,
      maxLines: maxLines,
      overflow: overflow,
    );
  }
}
`
}
