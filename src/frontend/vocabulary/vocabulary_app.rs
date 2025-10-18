use eframe::egui;

#[derive(Debug, Clone)]
pub struct VocabularyWord {
    pub word: String,
    pub translation: String,
    pub phonetic: String,
    pub example_sentence: String,
    pub proficiency: f32, // 熟练度 0.0 - 1.0
    pub last_reviewed: Option<std::time::SystemTime>,
}

pub struct VocabularyApp {
    words: Vec<VocabularyWord>,
    current_word_index: usize,
    show_translation: bool,
    new_word_input: String,
    new_translation_input: String,
    new_phonetic_input: String,
    new_example_input: String,
}

impl Default for VocabularyApp {
    fn default() -> Self {
        let words = vec![
            VocabularyWord {
                word: "apple".to_string(),
                translation: "苹果".to_string(),
                phonetic: "/ˈæp.əl/".to_string(),
                example_sentence: "I eat an apple every day.".to_string(),
                proficiency: 0.3,
                last_reviewed: None,
            },
            VocabularyWord {
                word: "book".to_string(),
                translation: "书".to_string(),
                phonetic: "/bʊk/".to_string(),
                example_sentence: "This is a good book.".to_string(),
                proficiency: 0.7,
                last_reviewed: None,
            },
            VocabularyWord {
                word: "computer".to_string(),
                translation: "计算机".to_string(),
                phonetic: "/kəmˈpjuː.tər/".to_string(),
                example_sentence: "I use a computer to work.".to_string(),
                proficiency: 0.5,
                last_reviewed: None,
            },
        ];

        Self {
            words,
            current_word_index: 0,
            show_translation: false,
            new_word_input: String::new(),
            new_translation_input: String::new(),
            new_phonetic_input: String::new(),
            new_example_input: String::new(),
        }
    }
}

impl VocabularyApp {
    pub fn show(&mut self, ui: &mut egui::Ui, back_to_main: &mut bool) {
        ui.heading("📚 背单词");
        ui.add_space(10.0);

        // 显示当前单词
        if !self.words.is_empty() {
            self.show_current_word(ui);
            ui.add_space(20.0);
            self.show_controls(ui);
        } else {
            ui.label("暂无单词，请添加新单词");
        }

        ui.add_space(20.0);

        // 显示添加新单词的表单
        self.show_add_word_form(ui);

        ui.add_space(20.0);

        // 显示单词列表
        self.show_word_list(ui);

        ui.add_space(20.0);

        // 返回按钮
        if ui.button("🔙 返回主菜单").clicked() {
            *back_to_main = true;
        }
    }

    fn show_current_word(&self, ui: &mut egui::Ui) {
        if self.words.is_empty() {
            return;
        }

        let word = &self.words[self.current_word_index];
        ui.heading(word.word.clone());
        ui.label(format!("音标: {}", word.phonetic));

        if self.show_translation {
            ui.add_space(10.0);
            ui.label(format!("翻译: {}", word.translation));
            ui.label(format!("例句: {}", word.example_sentence));
            ui.label(format!("熟练度: {:.1}%", word.proficiency * 100.0));
        }
    }

    fn show_controls(&mut self, ui: &mut egui::Ui) {
        ui.horizontal(|ui| {
            if ui.button("显示翻译").clicked() {
                self.show_translation = !self.show_translation;
            }

            if ui.button("认识").clicked() {
                if !self.words.is_empty() {
                    self.words[self.current_word_index].proficiency =
                        (self.words[self.current_word_index].proficiency + 0.1).min(1.0);
                    self.next_word();
                }
            }

            if ui.button("不认识").clicked() {
                if !self.words.is_empty() {
                    self.words[self.current_word_index].proficiency =
                        (self.words[self.current_word_index].proficiency - 0.1).max(0.0);
                    self.next_word();
                }
            }

            if ui.button("下一个").clicked() {
                self.next_word();
            }
        });
    }

    fn next_word(&mut self) {
        if !self.words.is_empty() {
            self.current_word_index = (self.current_word_index + 1) % self.words.len();
            self.show_translation = false;
        }
    }

    fn show_add_word_form(&mut self, ui: &mut egui::Ui) {
        ui.separator();
        ui.heading("➕ 添加新单词");

        ui.horizontal(|ui| {
            ui.label("单词:");
            ui.text_edit_singleline(&mut self.new_word_input);
        });
        ui.horizontal(|ui| {
            ui.label("翻译:");
            ui.text_edit_singleline(&mut self.new_translation_input);
        });
        ui.horizontal(|ui| {
            ui.label("音标:");
            ui.text_edit_singleline(&mut self.new_phonetic_input);
        });
        ui.horizontal(|ui| {
            ui.label("例句:");
            ui.text_edit_singleline(&mut self.new_example_input);
        });
        if ui.button("添加单词").clicked() {
            if !self.new_word_input.is_empty() && !self.new_translation_input.is_empty() {
                self.words.push(VocabularyWord {
                    word: self.new_word_input.clone(),
                    translation: self.new_translation_input.clone(),
                    phonetic: self.new_phonetic_input.clone(),
                    example_sentence: self.new_example_input.clone(),
                    proficiency: 0.0,
                    last_reviewed: None,
                });
                // 清空输入框
                self.new_word_input.clear();
                self.new_translation_input.clear();
                self.new_phonetic_input.clear();
                self.new_example_input.clear();
            }
        }
    }

    fn show_word_list(&self, ui: &mut egui::Ui) {
        ui.separator();
        ui.heading("📋 单词列表");
        egui::ScrollArea::vertical()
            .max_height(200.0)
            .show(ui, |ui| {
                for (i, word) in self.words.iter().enumerate() {
                    ui.horizontal(|ui| {
                        ui.label(format!("{}. {}", i + 1, word.word));
                        ui.label(format!("({})", word.translation));
                        ui.label(format!("[{:.0}%]", word.proficiency * 100.0));
                    });
                }
            });
    }
}
