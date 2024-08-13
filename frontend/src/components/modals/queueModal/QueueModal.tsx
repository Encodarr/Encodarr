import InputContainer from "../../inputs/inputContainer/InputContainer";
import Modal from "../../modal/Modal";

const QueueModal = ({ isOpen, setIsOpen, content, setContent }: any) => {
  const onClose = () => {
    setIsOpen(false);
  };
  const onModalSave = async () => {
    for (const key in content) {
      fetch(`/api/settings`, {
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
      onSave={onModalSave}
      title={"Queue Options"}
      onClose={onClose}
    >
      <InputContainer
        label="Page Size"
        type="text"
        selected={content.queue_page_size.value}
        onChange={(e: any) => {
          setContent({
            ...content,
            queue_page_size: e.target.value,
          });
        }}
      />
    </Modal>
  );
};
export default QueueModal;
