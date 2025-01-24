import InputContainer from "../../inputs/inputContainer/InputContainer";
import Modal from "../../modal/Modal";

const EventsModal = ({ isOpen, setIsOpen, content, setContent }: any) => {
	const onClose = () => {
		setIsOpen(false);
	};
	const onSave = async () => {
		for (const key in content) {
			fetch(`/api/settings/${key}`, {
				method: "PUT",
				headers: {
					"Content-Type": "application/json",
					Authorization: `Bearer ${localStorage.getItem("token")}`,
				},
				body: JSON.stringify({ id: key, value: content[key] }),
			});
		}
		setIsOpen(false);
	};
	if (!isOpen) return null;
	return (
		<Modal
			isOpen={isOpen}
			setIsOpen={setIsOpen}
			title={"Events Options"}
			onClose={onClose}
			onSave={onSave}
		>
			<InputContainer
				label="Page Size"
				type="text"
				selected={content.eventsPageSize}
				onChange={(e: any) => {
					setContent({
						...content,
						eventsPageSize: e.target.value,
					});
				}}
			/>
		</Modal>
	);
};
export default EventsModal;
