import styles from "./Events.module.scss";
import { SSEContext } from "../../contexts/webSocketContext";
import { useContext, useState } from "react";
import EventsModal from "../modals/eventsModal/EventsModal";
import EventsToolbar from "../toolbars/eventsToolbar/EventsToolbar";
import EventsTable from "../tables/eventsTable/EventsTable";

const Events = () => {
  const wsContext: any = useContext(SSEContext);
  const settings: any = wsContext?.data?.settings
    ? Object.keys(wsContext?.data?.settings).reduce((acc, key) => {
        acc[key] = wsContext?.data?.settings[key].value;
        return acc;
      }, {})
    : {};

  const logs = wsContext?.data?.logs || [];
  let sortedLogs = [...logs].sort((a, b) => b.id - a.id);
  if (settings?.eventsFilter == "info") {
    sortedLogs = sortedLogs.filter((log: any) => {
      return log.level == "INFO";
    });
  } else if (settings?.eventsFilter == "warn") {
    sortedLogs = sortedLogs.filter((log: any) => {
      return log.level == "WARN";
    });
  } else if (settings?.eventsFilter == "error") {
    sortedLogs = sortedLogs.filter((log: any) => {
      return log.level == "ERROR";
    });
  }
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [content, setContent] = useState<any>({});
  const [selected, setSelected] = useState<string | null>(null);
  return (
    <div className={styles.logs}>
      <EventsToolbar
        selected={selected}
        setContent={setContent}
        setIsModalOpen={setIsModalOpen}
        setSelected={setSelected}
        settings={settings}
      />
      <EventsModal
        isOpen={isModalOpen}
        setIsOpen={setIsModalOpen}
        content={content}
        setContent={setContent}
      />
      <div className={styles.content}>
        {sortedLogs && sortedLogs.length !== 0 ? (
          <>
            <EventsTable sortedLogs={sortedLogs} settings={settings} />
          </>
        ) : (
          <>Events are empty</>
        )}
      </div>
    </div>
  );
};
export default Events;
