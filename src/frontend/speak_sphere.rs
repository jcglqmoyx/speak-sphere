use crate::frontend::vocabulary::vocabulary_app::VocabularyApp;

pub(crate) struct SpeakSphere {
    pub(crate) current_screen: Screen,
    vocabulary_app: VocabularyApp,
}

#[derive(Default, PartialEq)]
pub(crate) enum Screen {
    #[default]
    Main,
    Vocabulary,
    Audiobook,
}

impl SpeakSphere {
    pub fn new() -> Self {
        Self {
            current_screen: Screen::Main,
            vocabulary_app: VocabularyApp::default(),
        }
    }
}

impl eframe::App for SpeakSphere {
    fn update(&mut self, ctx: &egui::Context, _frame: &mut eframe::Frame) {
        match self.current_screen {
            Screen::Main => self.show_main_screen(ctx),
            Screen::Vocabulary => self.show_vocabulary_screen(ctx),
            Screen::Audiobook => self.show_audiobook_screen(ctx),
        }
    }
}

impl SpeakSphere {
    fn show_main_screen(&mut self, ctx: &egui::Context) {
        egui::CentralPanel::default().show(ctx, |ui| {
            ui.heading("Speak Sphere - 语言学习平台");
            ui.add_space(20.0);

            // 按钮容器 - 使用水平布局让按钮并排显示
            ui.horizontal(|ui| {
                ui.add_space(ui.available_width() * 0.1); // 左边距

                // 左侧：背单词按钮
                ui.vertical(|ui| {
                    let vocab_button = ui.add_sized(
                        [200.0, 120.0],
                        egui::Button::new(egui::RichText::new("📚 背单词").heading()),
                    );

                    if vocab_button.clicked() {
                        self.current_screen = Screen::Vocabulary;
                    }

                    ui.label("智能记忆单词，轻松掌握词汇");
                    ui.label("包含不背单词功能");
                });

                ui.add_space(40.0); // 按钮间距

                // 右侧：阅读有声书按钮
                ui.vertical(|ui| {
                    let audiobook_button = ui.add_sized(
                        [200.0, 120.0],
                        egui::Button::new(egui::RichText::new("🎧 阅读有声书").heading()),
                    );

                    if audiobook_button.clicked() {
                        self.current_screen = Screen::Audiobook;
                    }

                    ui.label("沉浸式阅读，边听边学");
                    ui.label("集成轻听英语功能");
                });

                ui.add_space(ui.available_width() * 0.1); // 右边距
            });

            // 底部返回按钮（如果在子页面）
            if self.current_screen != Screen::Main {
                ui.add_space(20.0);
                if ui.button("🔙 返回主菜单").clicked() {
                    self.current_screen = Screen::Main;
                }
            }
        });
    }

    fn show_vocabulary_screen(&mut self, ctx: &egui::Context) {
        egui::CentralPanel::default().show(ctx, |ui| {
            let mut back_to_main = false;
            self.vocabulary_app.show(ui, &mut back_to_main);

            if back_to_main {
                self.current_screen = Screen::Main;
            }
        });
    }

    fn show_audiobook_screen(&mut self, ctx: &egui::Context) {
        egui::CentralPanel::default().show(ctx, |ui| {
            ui.heading("🎧 阅读有声书模块");
            ui.add_space(20.0);

            ui.label("这里是阅读有声书功能区域");
            ui.label("将集成轻听英语的核心功能：");
            ui.label("• 有声书籍阅读");
            ui.label("• 同步字幕显示");
            ui.label("• 语音控制");
            ui.label("• 生词标记");

            ui.add_space(20.0);
            if ui.button("🔙 返回主菜单").clicked() {
                self.current_screen = Screen::Main;
            }
        });
    }
}
