import InputContainer from "../../inputs/inputContainer/InputContainer";
import Modal from "../../modal/Modal";

const HistoryModal = ({ isOpen, setIsOpen, content, setContent }: any) => {
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
	const onClose = () => {
		setIsOpen(false);
	};
	if (!isOpen) return null;
	return (
		<Modal
			isOpen={isOpen}
			setIsOpen={setIsOpen}
			onSave={onSave}
			title={"History Options"}
			onClose={onClose}
		>
			<InputContainer
				label="Page Size"
				type="text"
				selected={content.historyPageSize}
				onChange={(e: any) => {
					setContent({
						...content,
						historyPageSize: e.target.value,
					});
				}}
			/>
		</Modal>
	);
};
export default HistoryModal;
