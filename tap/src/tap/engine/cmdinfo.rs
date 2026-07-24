use crate::engine::structs::{CmdInfo};

impl<'a, H> CmdInfo<'a, H> {
    pub fn new(
        handler: H,
        required: Vec<&'a str>,
        optional: Vec<&'a str>,
        unlimited: bool,
    ) -> Self {
        Self {
            handler,
            required,
            optional,
            unlimited,
        }
    }
}
