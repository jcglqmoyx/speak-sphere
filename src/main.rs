mod config;
mod utils;

use crate::config::get_display_size;
use crate::utils::font::setup_fonts;
use eframe::egui;

#[derive(Default)]
struct SpeakSphere {
    text_input: String,
}

impl eframe::App for SpeakSphere {
    fn update(&mut self, ctx: &egui::Context, _frame: &mut eframe::Frame) {
        egui::CentralPanel::default().show(ctx, |ui| {
            // ui.heading("多语言");
            // ui.horizontal(|ui| {
            //     let response = ui.text_edit_singleline(&mut self.text_input);
            //     if response.changed() {
            //         println!("输入内容: {}", self.text_input);
            //     }
            // });
        });
    }
}

fn main() -> Result<(), eframe::Error> {
    let (w, h) = get_display_size();
    let options = eframe::NativeOptions {
        viewport: egui::ViewportBuilder::default()
            .with_inner_size([(w as f32) * 0.6, (h as f32) * 0.6])
            .with_title("Speak Sphere"),
        ..Default::default()
    };

    eframe::run_native(
        "Speak Sphere",
        options,
        Box::new(|cc| {
            setup_fonts(&cc.egui_ctx);
            Ok(Box::<SpeakSphere>::default())
        }),
    )
}
