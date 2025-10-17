/// 简单的字体多语言支持检测
fn is_multilingual_font(path: &std::path::Path) -> bool {
    let filename = path
        .file_name()
        .and_then(|f| f.to_str())
        .unwrap_or("")
        .to_lowercase();

    // 检测字体文件名中的多语言关键词
    let multilingual_keywords = [
        "chinese",
        "chinese",
        "pingfang",
        "heiti",
        "sans",
        "hiragino",
        "korean",
        "korea",
        "hangul",
        "gothic",
        "myungjo",
        "japanese",
        "japan",
        "mincho",
        "gothic",
        "unicode",
        "symbol",
        "cjk",
        "extended",
        "arial",
        "helvetica",
        "tahoma",
        "verdana",
        "segoe",
    ];

    for keyword in &multilingual_keywords {
        if filename.contains(keyword) {
            return true;
        }
    }

    // 也包含一些知名字体文件名的直接匹配
    let known_fonts = [
        "arial unicode",
        "apple symbols",
        "pingfang",
        "stheit",
        "hiragino",
        "applegothic",
        "applemyungjo",
    ];

    for font_name in &known_fonts {
        if filename.contains(font_name) {
            return true;
        }
    }

    false
}

/// 加载一些通用的后备字体
fn load_fallback_fonts(fonts: &mut egui::FontDefinitions, _loaded_fonts: &mut Vec<String>) {
    // 这里可以添加一些在线字体下载逻辑，或者内置字体数据

    // 在实际应用中，可以考虑：
    // 1. 下载 Noto Sans 字体系列 (Google的开源字体)
    // 2. 使用内置的字体数据
    // 3. 提示用户安装字体

    // 使用_loaded_fonts参数避免警告
    let _ = _loaded_fonts;
    let _ = fonts;
}
/// 根据操作系统获取字体目录
fn get_system_font_directories() -> Vec<std::path::PathBuf> {
    let mut dirs = Vec::new();

    #[cfg(target_os = "macos")]
    {
        dirs.push(std::path::PathBuf::from("/System/Library/Fonts/"));
        dirs.push(std::path::PathBuf::from("/Library/Fonts/"));
        if let Some(home) = std::env::var_os("HOME") {
            let mut user_fonts = std::path::PathBuf::from(home);
            user_fonts.push("Library/Fonts/");
            dirs.push(user_fonts);
        }
    }

    #[cfg(target_os = "windows")]
    {
        if let Some(windir) = std::env::var_os("WINDIR") {
            let mut windows_fonts = std::path::PathBuf::from(windir);
            windows_fonts.push("Fonts");
            dirs.push(windows_fonts);
        }
    }

    #[cfg(target_os = "linux")]
    {
        dirs.push(std::path::PathBuf::from("/usr/share/fonts/"));
        dirs.push(std::path::PathBuf::from("/usr/local/share/fonts/"));
        if let Some(home) = std::env::var_os("HOME") {
            let mut user_fonts = std::path::PathBuf::from(home);
            user_fonts.push(".local/share/fonts/");
            dirs.push(user_fonts);
            user_fonts.pop();
            user_fonts.push(".fonts/");
            dirs.push(user_fonts);
        }
    }

    dirs
}
pub(crate) fn setup_fonts(ctx: &egui::Context) {
    // 获取默认字体定义
    let mut fonts = egui::FontDefinitions::default();

    // 根据操作系统获取字体目录
    let font_dirs = get_system_font_directories();
    let mut loaded_fonts = Vec::new();
    let mut font_index = 0;

    println!("扫描系统字体目录...");
    for font_dir in &font_dirs {
        println!("扫描目录: {}", font_dir.display());

        if let Ok(entries) = std::fs::read_dir(font_dir) {
            for entry in entries.flatten() {
                let path = entry.path();
                if let Some(extension) = path.extension() {
                    let ext_str = extension.to_string_lossy().to_lowercase();

                    // 支持的字体文件扩展名
                    if ext_str == "ttf" || ext_str == "ttc" || ext_str == "otf" {
                        if is_multilingual_font(&path) {
                            match std::fs::read(&path) {
                                Ok(font_bytes) => {
                                    let font_name = format!("system_font_{}", font_index);
                                    fonts.font_data.insert(
                                        font_name.clone(),
                                        egui::FontData::from_owned(font_bytes).into(),
                                    );
                                    loaded_fonts.push(font_name);
                                    println!("加载字体: {}", path.display());
                                    font_index += 1;

                                    // 限制加载的字体数量，避免内存占用过大
                                    if font_index >= 20 {
                                        break;
                                    }
                                }
                                Err(e) => {
                                    println!("读取字体失败 {}: {}", path.display(), e);
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    // 如果有成功加载的字体，将它们添加到字体家族中
    let font_count = loaded_fonts.len();
    if font_count > 0 {
        for font_name in &loaded_fonts {
            fonts
                .families
                .entry(egui::FontFamily::Proportional)
                .or_default()
                .insert(0, font_name.clone());
        }
        println!("✅ 已自动加载 {} 个系统字体", font_count);
    } else {
        println!("⚠️  未找到合适的系统字体，使用默认字体配置");

        // 尝试加载一些通用的后备字体路径
        load_fallback_fonts(&mut fonts, &mut loaded_fonts);
    }

    ctx.set_fonts(fonts);
}
