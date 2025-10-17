use rdev::display_size;

use once_cell::sync::OnceCell;

static DISPLAY_SIZE: OnceCell<(f32, f32)> = OnceCell::new();

pub(crate) fn get_unit_size() -> (f32, f32) {
    *DISPLAY_SIZE.get_or_init(|| {
        let (w, h) = display_size().expect("Failed to get display size");
        (w as f32 / 100.0, h as f32 / 100.0)
    })
}
