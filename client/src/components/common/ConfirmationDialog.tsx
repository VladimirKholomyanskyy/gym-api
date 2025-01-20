import { IconButton } from "@chakra-ui/react";
import {
  DialogActionTrigger,
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogRoot,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { FaTrash } from "react-icons/fa";
import { Button } from "../ui/button";

interface ConfirmationDialogProps {
  message: string;
  onDelete: () => void;
  title?: string;
  triggerLabel?: string;
  triggerIcon?: React.ReactNode;
  cancelLabel?: string;
  deleteLabel?: string;
}

const ConfirmationDialog = ({
  message,
  onDelete,
  title = "Are you sure?",
  triggerLabel = "Delete",
  triggerIcon = <FaTrash />,
  cancelLabel = "Cancel",
  deleteLabel = "Delete",
}: ConfirmationDialogProps) => {
  return (
    <DialogRoot role="alertdialog">
      <DialogTrigger asChild>
        {triggerIcon ? (
          <IconButton colorScheme="red" aria-label={triggerLabel}>
            {triggerIcon} {triggerLabel}
          </IconButton>
        ) : (
          <Button>{triggerLabel}</Button>
        )}
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
        </DialogHeader>
        <DialogBody>
          <p>{message}</p>
        </DialogBody>
        <DialogFooter>
          <DialogActionTrigger asChild>
            <Button variant="outline">{cancelLabel}</Button>
          </DialogActionTrigger>
          <DialogActionTrigger asChild>
            <Button colorScheme="red" onClick={onDelete}>
              {deleteLabel}
            </Button>
          </DialogActionTrigger>
        </DialogFooter>
        <DialogCloseTrigger />
      </DialogContent>
    </DialogRoot>
  );
};

export default ConfirmationDialog;
