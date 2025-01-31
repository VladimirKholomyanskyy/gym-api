import { IconButton, Stack, Text } from "@chakra-ui/react";
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

interface DeleteDialogProps {
  message: string;
  onDelete: () => void;
  title?: string;
}

const DeleteDialog = ({
  message,
  onDelete,
  title = "Are you sure?",
}: DeleteDialogProps) => {
  return (
    <DialogRoot role="alertdialog">
      <DialogTrigger asChild>
        <IconButton
          background="transparent"
          color="red.600"
          _hover={{ color: "neon.300" }}
          aria-label="Delete"
        >
          <FaTrash /> Delete
        </IconButton>
      </DialogTrigger>
      <DialogContent
        background="blackAlpha.900"
        border="1px solid"
        borderColor="neon.400"
        boxShadow="0 0 15px rgba(0, 255, 255, 0.8)"
        p={4}
      >
        <DialogHeader>
          <DialogTitle
            color="neon.400"
            textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
          >
            {title}
          </DialogTitle>
        </DialogHeader>
        <DialogBody>
          <Text color="gray.300">{message}</Text>
        </DialogBody>
        <DialogFooter>
          <Stack direction="row" gap={4}>
            <DialogActionTrigger asChild>
              <Button
                variant="outline"
                borderColor="neon.400"
                color="neon.400"
                _hover={{ borderColor: "neon.300", color: "neon.300" }}
              >
                Cancel
              </Button>
            </DialogActionTrigger>
            <DialogActionTrigger asChild>
              <Button
                background="red.600"
                color="white"
                _hover={{ background: "red.400", boxShadow: "0 0 15px red" }}
                onClick={onDelete}
              >
                Delete
              </Button>
            </DialogActionTrigger>
          </Stack>
        </DialogFooter>
        <DialogCloseTrigger />
      </DialogContent>
    </DialogRoot>
  );
};

export default DeleteDialog;
